package mohttp

import (
	"fmt"
)

type HTTPError struct {
	Code   int
	Reason string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error code=%d %s", e.Code, e.Reason)
}
