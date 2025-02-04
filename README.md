1) Trabajar en ubuntu 22.04 LTS

2) Tener instalado Docker: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04

3) Clonar repositorio

4) Ir a carpeta

4) Crear archivo variables de entorno

5) dar permiso a carpeta data: chmod a+rwx ./data

6) Descargar datos (Se demora entre 5 y 10 minutos)

7) Ejecutar con docker compose up --build

 ver proyectos en su respectivos host

8) Ejecutar indexer (sin necesidad de go)

Nota: Analisar si paso 5, 6 y 8 se pueden poner en un archivo.sh y ejecutar desde docker compose, teniendo en cuenta que para que suceda la ejecucion de la base de datos ya se debe dar permiso a la carpeta data (es decir este contenedor nuevo si lo es deberia ser una dependencia del contenedor zincsearch) (Hay un problema si se crea un contenedor todavia no se ha descargado los datos y no sabemos si eso sea un problema es posible que tenga que ver con volumenes), luego para crear el indexer ya se debe ejecutar el descargador de datos y la base de datos ya debe estar instalada
Instalar go
https://wiki.crowncloud.net/?How_To_Install_Go_on_Ubuntu_24_04

go build -o indexer ./indexer/indexer.go

#Analizar profiling


go tool pprof -top cpu_profile.prof
go tool pprof cpu_profile.prof (No es necesario)
go tool pprof -svg cpu_profile.prof > cpu_profile.svg
go tool pprof -top mem_profile.prof
go tool pprof -svg mem_profile.prof > mem_profile.svg
go tool pprof -http=:8090 cpu_profile.prof

