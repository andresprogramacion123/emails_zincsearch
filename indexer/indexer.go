package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"runtime/pprof" // 🛠️ Importar pprof para profiling
	"bufio"
)

// 📌 Archivo donde se guardará el profiling
var cpuProfile = "./indexer/cpu_profile.prof"
var memProfile = "./indexer/mem_profile.prof"

type EmailData struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	MessageID   string `json:"message_id"`
	Date        string `json:"date"`
	ContentType string `json:"content_type"`
	OfFolder    string `json:"of_folder"`
}

type BulkData struct {
	Index   string      `json:"index"`
	Records []EmailData `json:"records"`
}

type PropertyDetail struct {
	Type          string `json:"type"`
	Index         bool   `json:"index"`
	Store         bool   `json:"store"`
	Sortable      bool   `json:"sortable"`
	Aggregatable  bool   `json:"aggregatable"`
	Highlightable bool   `json:"highlightable"`
}

type Mapping struct {
	Properties map[string]PropertyDetail `json:"properties"`
}

type IndexerData struct {
	Name         string  `json:"name"`
	StorageType  string  `json:"storage_type"`
	ShardNum     int     `json:"shard_num"`
	MappingField Mapping `json:"mappings"`
}

// Cargar variables de entorno desde un archivo .env
func loadEnvFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("No se pudo abrir el archivo .env: %v\n", err)
		return // No detener el programa si el .env no existe
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías o comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Dividir en clave y valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Línea malformada en .env: %s\n", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remover comillas si existen
		value = strings.Trim(value, `"'`)

		// Establecer la variable de entorno
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error leyendo el archivo .env: %v\n", err)
	}
}

func main() {
	// Cargar variables de entorno desde el archivo .env
	loadEnvFile(".env")

	// 🛠️ Habilitar profiling de CPU
	f, err := os.Create(cpuProfile)
	if err != nil {
		log.Fatal("No se pudo crear el archivo de profiling de CPU:", err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// 🛠️ Habilitar profiling de memoria al finalizar
	defer func() {
		memFile, err := os.Create(memProfile)
		if err != nil {
			log.Fatal("No se pudo crear el archivo de profiling de memoria:", err)
		}
		defer memFile.Close()
		pprof.WriteHeapProfile(memFile)
	}()

	log.Println("Starting indexer!")
	//Ingesamos nombre de nuestro indice json en proyecto
	indexerData, err := createIndexerFromJsonFile("./indexer/julian_indexer.json")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Deleting index if exists...")
	//Ingesamos nombre de nuestro indice en zincsearch para ver si existe y ser eliminado
	deleted := deleteIndexOnZincSearch("julian_emails")
	if deleted != nil {
		log.Println("Index doesn't exist. Creating...")
	}

	sent := createIndexOnZincSearch(indexerData)
	if sent != nil {
		log.Fatal(sent)
	}

	log.Println("Index created successfully.")
	log.Println("Start indexing, this might take a few minutes...")
	startTime := time.Now()

	var records []EmailData
	var m sync.Mutex
	var wg sync.WaitGroup

	//Ingresamos ruta de datos de emails
	err = filepath.Walk("./indexer/enron_mail_20110402/maildir/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				emailData, err := processFile(p)
				
				if err != nil {
					log.Println(err)
					return
				}
				m.Lock()
				records = append(records, emailData)
				m.Unlock()
			}(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	
	// Define el tamaño del lote que funciona bien en tu máquina
	const batchSize = 5000 
	sendBulkToZincSearch(records, batchSize)

	duration := time.Since(startTime)
	log.Printf("Finished indexing. Time taken: %.2f seconds", duration.Seconds())
}

func createIndexerFromJsonFile(filepath string) (IndexerData, error) {
	var indexerData IndexerData

	file, err := os.Open(filepath)
	if err != nil {
		return indexerData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&indexerData)
	if err != nil {
		return indexerData, err
	}

	return indexerData, nil
}

func deleteIndexOnZincSearch(indexName string) error {
	zincHost := os.Getenv("HOST")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:4080/api/index/%s", zincHost, indexName), nil)
	if err != nil {
		return err
	}
	
	zincUser := os.Getenv("ZINC_USER")
	zincPassword := os.Getenv("ZINC_PASSWORD")
	//Ingresamos usuario y contraseña de db zincsearch
	req.SetBasicAuth(zincUser, zincPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete indexer, status code: %d", resp.StatusCode)
	}

	log.Println("Index deleted successfully")
	return nil
}

func createIndexOnZincSearch(indexerData IndexerData) error {
	jsonData, err := json.Marshal(indexerData)
	if err != nil {
		log.Fatal(err)
	}
	zincHost := os.Getenv("HOST")
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:4080/api/index", zincHost), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	//Ingresamos usuario y contraseña de db zincsearch
	zincUser := os.Getenv("ZINC_USER")
	zincPassword := os.Getenv("ZINC_PASSWORD")
	req.SetBasicAuth(zincUser, zincPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to create indexer, status code: %d", resp.StatusCode)
	}

	return nil
}

func processFile(path string) (EmailData, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return EmailData{}, fmt.Errorf("Error leyendo archivo %s: %v", path, err)
	}

	// 🔹 Extraer el usuario desde la ruta del archivo
	userFolder := extractUserFolder(path)

	// 🔹 Limpiar encabezados mal formateados
	cleanedContent := cleanMalformedHeaders(string(content))

	msg, err := mail.ReadMessage(bytes.NewReader([]byte(cleanedContent)))
	if err != nil {
		log.Printf("Error procesando archivo %s: %v", path, err)
		return EmailData{}, err
	}

	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")
	messageID := msg.Header.Get("Message-ID")
	date := msg.Header.Get("Date")
	contentType := msg.Header.Get("Content-Type")

	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("Error leyendo el cuerpo del mensaje en %s: %v", path, err)
		return EmailData{}, err
	}

	return EmailData{
		From:        from,
		To:          to,
		Subject:     subject,
		Content:     strings.TrimSpace(string(body)),
		MessageID:   messageID,
		Date:        date,
		ContentType: contentType,
		OfFolder:    userFolder, // 🔹 Agregar la carpeta del usuario
	}, nil
}

// 🔹 Función para extraer la carpeta del usuario desde la ruta
func extractUserFolder(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	for i, part := range parts {
		if part == "maildir" && i+1 < len(parts) {
			return parts[i+1] // 🔹 Retorna la carpeta inmediatamente después de "maildir"
		}
	}
	return "unknown" // Si no se encuentra, retorna "unknown"
}

func cleanMalformedHeaders(content string) string {
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	inHeaders := true
	prevLineIndex := -1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 🔹 Si es una línea vacía, marca el fin de los encabezados
		if trimmed == "" {
			inHeaders = false
		}

		if inHeaders {
			// 🔹 Si la línea NO tiene `:` y NO es la primera línea del email, es parte del encabezado anterior
			if !strings.Contains(trimmed, ":") && prevLineIndex != -1 {
				cleanedLines[prevLineIndex] += " " + trimmed // Agrega el contenido a la línea anterior
				continue
			}

			// 🔹 Guardamos el índice de la última línea válida con `:`
			if strings.Contains(trimmed, ":") {
				prevLineIndex = len(cleanedLines)
			}
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}

func sendBulkToZincSearch(records []EmailData, batchSize int) {
    totalRecords := len(records)
    log.Printf("Total records to index: %d\n", totalRecords)
	
    for start := 0; start < totalRecords; start += batchSize {
        end := start + batchSize
        if end > totalRecords {
            end = totalRecords
        }

        batch := records[start:end] // Tomamos el lote de registros
        bulkData := BulkData{
            Index:   "julian_emails", //Ingresamos nombre de indice en zincsearch (poner como variable variable env)
            Records: batch,
        }

        jsonData, err := json.Marshal(bulkData)
        if err != nil {
            log.Println(err)
            return
        }

		zincHost := os.Getenv("HOST")
        req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:4080/api/_bulkv2", zincHost), bytes.NewReader(jsonData))
        if err != nil {
            log.Println(err)
            return
        }
		
		//Ingresamos usuario y contraseña de db zincsearch
		zincUser := os.Getenv("ZINC_USER")
		zincPassword := os.Getenv("ZINC_PASSWORD")
        req.SetBasicAuth(zincUser, zincPassword)
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Println(err)
            return
        }
		
        defer resp.Body.Close()

        _, err = io.ReadAll(resp.Body)
        if err != nil {
            log.Println(err)
            return
        }
		
        log.Printf("Batch from %d to %d indexed successfully\n", start, end)

        // **Liberar memoria del lote procesado**
        batch = nil
    }
}





