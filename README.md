1) Trabajar en ubuntu 22.04 LTS

2) Tener instalado Docker: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04

3) Clonar repositorio

4) Ir a carpeta

4) Crear archivo variables de entorno

5) dar permiso a carpeta data: chmod a+rwx ./data

6) Descargar datos (Se demora entre 5 y 10 minutos) (para desarrollo)

7) Ejecutar con docker compose up --build

 ver proyectos en su respectivos host

8) Ejecutar indexer (sin necesidad de go)

./indexer/indexer

Con go

Instalar go
https://wiki.crowncloud.net/?How_To_Install_Go_on_Ubuntu_24_04

go build -o indexer ./indexer/indexer.go

Despues de que los datos estan indexados podemos finalizar la ejecucion del servicio de zincsearch.

Generar un backup

Este backup es necesario para produccion

#Analizar profiling

go tool pprof -top cpu_profile.prof
go tool pprof cpu_profile.prof (No es necesario)
go tool pprof -svg cpu_profile.prof > cpu_profile.svg
go tool pprof -top mem_profile.prof
go tool pprof -svg mem_profile.prof > mem_profile.svg
go tool pprof -http=:8090 cpu_profile.prof

Despliegue:

Conexion ssh
ssh -i clave-julian.pem ubuntu@3.83.80.213

Instalar docker

Clonar repo

dar permiso a data

crear variables de entorno

Crear .env en cliente y cambiar fecth en cliente

Copiar backup en data
scp -i clave-julian.pem -r backup_data/ ubuntu@52.91.213.148:/home/ubuntu/

