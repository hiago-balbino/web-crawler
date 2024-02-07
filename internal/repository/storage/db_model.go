package storage

type pageDataInfo struct {
	URI   string   `bson:"uri"`
	Depth uint     `bson:"depth"`
	URIs  []string `bson:"uris"`
}
