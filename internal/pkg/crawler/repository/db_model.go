package crawler

type dataPage struct {
	URI   string   `bson:"uri"`
	Depth uint     `bson:"depth"`
	URIs  []string `bson:"uris"`
}
