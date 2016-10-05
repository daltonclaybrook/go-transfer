package middle

import (
	"net/http"
)

// CORS is a piece of middleware used to add CORS headers to responses
type CORS struct {
	Origin  string
	Methods string
	Headers string
}

// Handle adds CORS headers to a response
func (c CORS) Handle(w http.ResponseWriter, r *http.Request, context Context, next NextFunc) {
	w.Header().Set("Access-Control-Allow-Origin", c.Origin)
	w.Header().Set("Access-Control-Allow-Methods", c.Methods)
	w.Header().Set("Access-Control-Allow-Headers", "content-type, content-length")
	next(context)
}
