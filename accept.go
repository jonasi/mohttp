package mohttp

import (
	"golang.org/x/net/context"
	"mime"
	"strings"
)

type acceptDataMapper struct {
}

func (a *acceptDataMapper) HandleErr(c context.Context, err error) {

}

func (a *acceptDataMapper) HandleResult(c context.Context, res interface{}) error {
	return nil
}

func AcceptDataMapper(m map[string]DataResponder) Handler {
	am := &acceptDataMapper{}

	return HandlerFunc(func(c context.Context) {
		c = WithResponder(c, am)
		Next(c)
	})
}

func RequireMediaTypeHandler(mediaTypes ...string) Handler {
	mediaTypes = ParseMediaTypes(mediaTypes...)

	return HandlerFunc(func(c context.Context) {
		if !MatchMediaTypes(c, mediaTypes...) {
			Error(c, "Not Acceptable", 406)
			return
		}
	})
}

func AcceptHandler(mediaTypes ...string) func(Handler) Handler {
	mediaTypes = ParseMediaTypes(mediaTypes...)

	return func(h Handler) Handler {
		return HandlerFunc(func(c context.Context) {
			if MatchMediaTypes(c, mediaTypes...) {
				h.Handle(c)
				return
			}

			Error(c, "Not Acceptable", 406)
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
	acc := strings.Split(GetRequest(c).Header.Get("Accept"), ",")

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
