package internal

import (
	"fmt"
	"os"
)

func WrapMain(mainFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(os.Stderr, "an error occured: %s", err)
			} else {
				fmt.Fprintf(os.Stderr, "something unexpected happened (%T), %s", r, r)
			}
			os.Exit(1)
		}
	}()
	mainFunc()
}
