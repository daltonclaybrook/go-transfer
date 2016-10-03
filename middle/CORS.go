package middle

import (
	"net/http"
)

type CORS struct {
	Origin  string
	Methods string
}

func (c CORS) Handle(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.Origin)
		w.Header().Set("Access-Control-Allow-Methods", c.Methods)
		handler(w, r)
	}
}
