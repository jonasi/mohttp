package middleware

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"mime"
	"strings"
)

type AcceptHandlers map[string]mohttp.Handler

func (a AcceptHandlers) Handle(c context.Context) {
	for mt, h := range a {
		if MatchMediaTypes(c, mt) {
			h.Handle(c)
			mohttp.Next(c)
			return
		}
	}

	mohttp.Error(c, "Not Acceptable", 406)
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
