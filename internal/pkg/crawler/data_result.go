package crawler

type dataResult struct {
	uri   string
	uris  []string
	err   error
	depth int32
}
