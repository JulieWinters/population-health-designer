package maps

import (
	"errors"
	"fmt"
)

func Execute(config string) (string, error) {
	return fmt.Sprintf("ECHO '%v'", config), errors.New("Unimplemented")
}
