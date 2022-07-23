package binding

import (
	"net/http"
)

type QueryBinding struct{}

func (QueryBinding) Name() string {
	return "query"
}

func (QueryBinding) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	return decoder.Decode(obj, values)
}
