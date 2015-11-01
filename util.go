package mohttp

import (
	"fmt"
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
