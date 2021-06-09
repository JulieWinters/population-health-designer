package maps

import (
	"errors"
	"fmt"
)

func Execute(config string) (string, error) {
	return fmt.Sprint("ECHO '%v'", config), errors.New("Unimplemented")
}
