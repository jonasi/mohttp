package http

import (
	"net/http"
)

func StripPrefix(prefix string) Handler {
	return FromStdLib(func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	})
}
