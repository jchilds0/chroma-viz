#!/bin/bash 

go build main.go 
go run main.go -cpuprofile=perf/gui.prof
go tool pprof -nodecount 10 -top main perf/gui.prof
go tool pprof -http=localhost:8080 main perf/gui.prof
