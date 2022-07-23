package binding

import (
	"net/http"
)

const defaultMemory = 32 << 20

type FormBinding struct{}

type FormPostBinding struct{}

type FormMultipartBinding struct{}

func (FormPostBinding) Name() string {
	return "form-urlencoded"
}

func (FormPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	return decoder.Decode(obj, req.PostForm)
}

func (FormMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (FormMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	return decoder.Decode(obj, req.PostForm)
}
