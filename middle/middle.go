package middle

import (
	"net/http"
)

type Middle interface {
	Handle(handler http.HandlerFunc) http.HandlerFunc
}
