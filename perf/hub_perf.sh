#!/bin/bash

go run ./cmd/chroma-hub -profile &
go tool pprof -http=localhost:8080 http://localhost:9000/debug/pprof/profile &
go run ./cmd/chroma-viz -t 100
