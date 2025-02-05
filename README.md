# Prueba tecnica desarrollador de software: Interfaz para buscar informacion de bases de datos de correos

Por: Julian Andres Montoya Carvajal

C.C 1214727927

## Estructura de proyecto:

## Instrucciones etapa de desarrollo: 
1) Preferiblemente trabajar en ubuntu 22.04 LTS (Puede utilizar wsl2 en windows para instalar ubuntu)

2) Tener instalado Docker y Docker compose para la orquestacion de los servicios: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04

3) Clone el repositorio con el siguiente comando:

```
git clone https://github.com/andresprogramacion123/emails_zincsearch.git
```

4) Ingrese a la carpeta del proyecto emails_zincsearch

4) Crear archivo variables de entorno .env en la raiz del proyecto para el caso del backend y base de datos zincsearch

```
touch .env
```

Escriba las variables de entorno en archivo .env correspondientes a ZINC_USER, ZINC_PASSWORD y HOST.

Para el caso del frontend cree el archivo .env dentro de la carpeta client y establezca la variable VITE_API_URL correspondiente a la URL de la API Rest.

5) Dar permiso a carpeta data. Esta es la carpeta correspondiente para almacenar datos y metadatos de zincsearch y el indixe.

```
chmod a+rwx ./data
```

6) Descargar datos de enron_email (se demora entre 5 y 10 minutos, de acuerdo las capacidades): Para descargar los datos desde la terminal ejecute la siguiente instruccion:

```
./indexer/download_data.sh
```

Esto descargara en la carpeta indexer los datos correspondientes a 517424 correos.

7) Una vez los datos se encuentran descargados, podemos ejecutar los tres servicios necesarios para la aplicacion: Base de datos Zincsearch, API Rest y Frontend haciendo uso de docker y docker compose. Para ello ejecute en la terminal el siguiente comando

```
sudo docker compose up --build
```

En caso de ser necesario ejecute antes los siguientes comandos:

```
sudo docker compose down
```

```
sudo docker compose down -v
```

```
sudo docker compose down --rmi all
```

```
sudo docker builder prune -a -f
```

A continuacion puede ver los servicios ejecutandose en los siguientes enlaces:

* Interfaz grafica ZincSearch: http://localhost:4080/
* API Rest: http://localhost:8080/
* Frontend: http://localhost:5173/

8) Con la base de datos ZincSearch ejecutandose podemos pasar a crear e indexer los datos (Puede tardar un minuto dependiendo de los recursos y del tamaño del lote).

* Si desea compilar el indice desde cero:

Debe tener instalado Go: https://wiki.crowncloud.net/?How_To_Install_Go_on_Ubuntu_24_04

Ejecutar el siguiente comando para crear el compilable llamado indexer

```
go build -o indexer ./indexer/indexer.go
```

Ejecutar el compilador:

```
./indexer/indexer

```

* Si solo desea ejecutar el indice ya establecido

```
./indexer/indexer
```

Teniendo los datos indexados podemos proceder a visualizarlos en el frontend o el la interfaz grafica de zincsearch de manera local

9) Procedemos a generar graficos de profiling:

```
go tool pprof -top ./indexer/cpu_profile.prof
```

```
go tool pprof -svg ./indexer/cpu_profile.prof > ./indexer/cpu_profile.svg
```

```
go tool pprof -top ./indexer/mem_profile.prof

```

```
go tool pprof -svg ./indexer/mem_profile.prof > ./indexer/mem_profile.svg
```

```
go tool pprof -http=:8090 ./indexer/cpu_profile.prof
```

```
go tool pprof -http=:8090 ./indexer/mem_profile.prof
```

10) Despues de que los datos estan indexados podemos finalizar la ejecucion del servicio de zincsearch.

11) Generar un backup: Este backup es necesario para produccion (Vamos aca)
Falta escribir...

## Instrucciones etapa de produccion:

Adquirir un servidor de ec2.

Conexion ssh
ssh -i clave-julian.pem ubuntu@52.91.213.148

Instalar docker

Clonar repo

dar permiso a data

crear variables de entorno

Crear .env en cliente y cambiar fecth en cliente

Copiar backup en data
scp -i clave-julian.pem -r backup_data/ ubuntu@52.91.213.148:/home/ubuntu/

