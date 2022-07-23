package binding

import (
	"github.com/gorilla/schema"
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

var decoder = schema.NewDecoder()

func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return QueryBinding{}
	}

	switch contentType {
	case MIMEJSON:
		return JSONBinding{}
	case MIMEPOSTForm:
		return FormPostBinding{}
	default:
		return FormMultipartBinding{}
	}
}

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}
