package hub

import "os"

type Asset struct {
	Filename  string
	Directory string
	Name      string
	ImageID   int
	image     []byte
}

func (asset *Asset) fetchImage() (err error) {
	asset.image, err = os.ReadFile(asset.Filename)
	return err
}

type Assets map[int]Asset
