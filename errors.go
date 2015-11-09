package mohttp

import (
	"fmt"
	"net/http"
)

type HTTPError int

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d %s", int(e), http.StatusText(int(e)))
}
