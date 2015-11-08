package mohttp

import (
	"fmt"
	"golang.org/x/net/context"
)

func recoverErr() error {
	r := recover()

	if r == nil {
		return nil
	}

	if err, ok := r.(error); ok {
		return err
	}

	return fmt.Errorf("%#v", r)
}

func RouteID(c context.Context) string {
	route, _ := GetRoute(c)
	return fmt.Sprintf("%s.%s.%s", RouterID(c), route.Method(), route.Path())
}

func RouterID(c context.Context) string {
	router, _ := GetRouter(c)
	return fmt.Sprintf("%p", router)
}
