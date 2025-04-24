package hub

import (
	"encoding/json"
	"fmt"
	"os"
)

type Asset struct {
	AssetID   int64
	Directory string
	Name      string
	Image     []byte
}

type AssetRef struct {
	AssetID   int
	Directory string
	Name      string
	FileName  string
}

func AssetsFromFile(fileName string) ([]Asset, error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	assetRefs := make([]AssetRef, 0, 100)
	err = json.Unmarshal(buf, &assetRefs)
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, 0, len(assetRefs))
	for _, ref := range assetRefs {
		asset := Asset{
			AssetID:   int64(ref.AssetID),
			Name:      ref.Name,
			Directory: ref.Directory,
		}

		asset.Image, err = os.ReadFile(ref.FileName)
		if err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (hub *DataBase) GetAssets() (assets []Asset, err error) {
	rows, err := hub.db.Query("SELECT a.assetID FROM asset a;")
	if err != nil {
		return
	}

	assets = make([]Asset, 0, 10)
	var assetID int64

	for rows.Next() {
		err = rows.Scan(&assetID)
		if err != nil {
			err = fmt.Errorf("AssetID: %s", err)
			return
		}

		asset, err := hub.GetAsset(assetID)
		if err != nil {
			err = fmt.Errorf("retrieve asset: %s", err)
		}

		asset.Image = []byte{}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (hub *DataBase) GetAsset(assetID int64) (asset Asset, err error) {
	hub.lock.Lock()
	asset, ok := hub.assets[assetID]
	hub.lock.Unlock()
	if ok {
		return
	}

	asset.AssetID = assetID

	row := hub.stmt[ASSET_SELECT].QueryRow(assetID)
	err = row.Scan(&asset.Name, &asset.Directory, &asset.Image)
	if err != nil {
		err = fmt.Errorf("Asset %d: %s", assetID, err)
		return
	}

	hub.lock.Lock()
	hub.assets[assetID] = asset
	hub.lock.Unlock()

	return
}

func (hub *DataBase) ImportAsset(asset Asset) (err error) {
	_, err = hub.stmt[ASSET_DELETE].Exec(asset.AssetID)
	if err != nil {
		Logger(err.Error())
	}

	_, err = hub.stmt[ASSET_INSERT].Exec(asset.AssetID, asset.Name, asset.Directory, asset.Image)
	return
}
