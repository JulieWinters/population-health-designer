package maps

import (
	"errors"
	"fmt"
)

// Executes the datd mapping process
func Execute(config string) (string, error) {
	return fmt.Sprintf("ECHO '%v'", config), errors.New("Unimplemented")
}
