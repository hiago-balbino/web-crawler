package handler

type crawPageInfo struct {
	URI   string `form:"uri"`
	Depth uint   `form:"depth"`
}

func (cp crawPageInfo) validate() error {
	switch {
	case cp.URI == "":
		return errEmptyURI
	case cp.Depth == 0:
		return errEmptyDepth
	default:
		return nil
	}
}
