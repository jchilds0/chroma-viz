.PHONY: artist
artist:
	go build -o ./bin ./cmd/artist/

.PHONY: hub
hub: 
	go build -o ./bin ./cmd/chroma_hub/

.PHONY: viz
viz: 
	go build -o ./bin ./cmd/viz/

all: artist hub viz
