.PHONY: artist hub viz build
artist:
	go build -o ./bin ./cmd/chroma-artist/

hub: 
	go build -o ./bin ./cmd/chroma-hub/

viz: 
	go build -o ./bin ./cmd/chroma-viz/

build: 
	go build -o ./bin ./cmd/chroma-viz/
	go build -o ./bin ./cmd/chroma-hub/
	go build -o ./bin ./cmd/chroma-artist/

