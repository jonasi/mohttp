package middleware

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"mime"
	"strings"
)

type AcceptDataMapper map[string]mohttp.Handler

func (a AcceptDataMapper) Handle(c context.Context) {

}

func (a AcceptDataMapper) HandleErr(c context.Context, err error) {

}

func (a AcceptDataMapper) HandleResult(c context.Context, res interface{}) error {
	return nil
}

func RequireMediaTypeHandler(mediaTypes ...string) mohttp.Handler {
	mediaTypes = ParseMediaTypes(mediaTypes...)

	return mohttp.HandlerFunc(func(c context.Context) {
		if !MatchMediaTypes(c, mediaTypes...) {
			mohttp.Error(c, "Not Acceptable", 406)
			return
		}
	})
}

func AcceptHandler(mediaTypes ...string) func(mohttp.Handler) mohttp.Handler {
	mediaTypes = ParseMediaTypes(mediaTypes...)

	return func(h mohttp.Handler) mohttp.Handler {
		return mohttp.HandlerFunc(func(c context.Context) {
			if MatchMediaTypes(c, mediaTypes...) {
				h.Handle(c)
				return
			}

			mohttp.Error(c, "Not Acceptable", 406)
		})
	}
}

func ParseMediaTypes(mediaTypes ...string) []string {
	for i, mt := range mediaTypes {
		t, _, err := mime.ParseMediaType(mt)

		if err != nil {
			panic("Media type parsing error: " + err.Error())
		}

		mediaTypes[i] = t
	}

	return mediaTypes
}

func MatchMediaTypes(c context.Context, mediaTypes ...string) bool {
	acc := strings.Split(mohttp.GetRequest(c).Header.Get("Accept"), ",")

	if len(acc) == 0 {
		return true
	}

	for _, mt := range mediaTypes {
		if mt == "*/*" {
			return true
		}

		for _, h := range acc {
			if mt == strings.TrimSpace(h) {
				return true
			}
		}
	}

	return false
}
