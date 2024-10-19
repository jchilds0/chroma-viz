#!/bin/bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

tmux neww -n ChromaHub 'go run ./cmd/chroma-hub -profile perf'
sleep 3
tmux neww -n pprof 'go tool pprof -http=127.0.0.1:8080 http://127.0.0.1:9000/debug/pprof/profile'
tmux neww -n ChromaViz 'go run ./cmd/chroma-viz -t 100'
