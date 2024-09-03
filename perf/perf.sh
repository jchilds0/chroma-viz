#!/bin/bash 

file=perf/viz.prof

go run ./cmd/chroma-viz -t 1000 -profile $file
go tool pprof -nodecount 10 -top main $file
go tool pprof -http=localhost:8080 main $file
