#!/bin/bash 

file=perf/viz.prof

go run ./cmd/chroma-viz -profile $file
go tool pprof -http=localhost:8080 ./cmd/chroma-viz $file
