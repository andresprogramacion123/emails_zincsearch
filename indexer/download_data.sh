#!/bin/bash

# Definir la URL y el directorio de destino
url="http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz"
destination_dir="$(dirname "$0")"  # Obtiene la ruta absoluta de la carpeta actual (indexer)
file_name="enron_mail_20110402.tgz"

# Descargar el archivo dentro de la carpeta indexer
echo "Descargando datos en $destination_dir..."
wget -O "$destination_dir/$file_name" "$url"

# Verificar si la descarga fue exitosa
if [ $? -eq 0 ]; then
    echo "Datos descargados correctamente en $destination_dir/$file_name"
else
    echo "Error al descargar los datos"
    exit 1
fi

# Descomprimir dentro de la carpeta indexer
echo "Descomprimiendo en $destination_dir..."
tar -xzvf "$destination_dir/$file_name" -C "$destination_dir"

echo "Datos descomprimidos en $destination_dir"
echo "Listo! Ahora puedes trabajar con tus datos."

