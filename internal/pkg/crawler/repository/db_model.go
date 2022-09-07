package crawler

type dataPage struct {
	URI   string   `bson:"uri"`
	Depth int      `bson:"depth"`
	URIs  []string `bson:"uris"`
}
