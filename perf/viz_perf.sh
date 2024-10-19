#!/bin/bash 

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

file=perf/viz.prof

tmux neww -n ChromaHub "go run ./cmd/chroma-hub"
sleep 3
go run ./cmd/chroma-viz -profile $file
go tool pprof -http=127.0.0.1:8080 $file
