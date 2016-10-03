package middle

import (
	"net/http"
)

type CORS struct {
	Origin  string
	Methods string
	Headers string
}

func (c CORS) Handle(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.Origin)
		w.Header().Set("Access-Control-Allow-Methods", c.Methods)
		w.Header().Set("Access-Control-Allow-Headers", "content-type, content-length")
		handler(w, r)
	}
}
