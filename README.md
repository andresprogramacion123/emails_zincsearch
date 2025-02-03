go build -o indexer ./indexer/indexer.go

#Analizar profiling


go tool pprof -top cpu_profile.prof
go tool pprof cpu_profile.prof (No es necesario)
go tool pprof -svg cpu_profile.prof > cpu_profile.svg
go tool pprof -top mem_profile.prof
go tool pprof -svg mem_profile.prof > mem_profile.svg
go tool pprof -http=:8090 cpu_profile.prof

