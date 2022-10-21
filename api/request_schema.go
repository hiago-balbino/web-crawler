package api

type requestSchema struct {
	URI   string `form:"uri"`
	Depth uint   `form:"depth"`
}

func (rs requestSchema) validate() error {
	switch {
	case rs.URI == "":
		return errEmptyURI
	case rs.Depth == 0:
		return errEmptyDepth
	default:
		return nil
	}
}
