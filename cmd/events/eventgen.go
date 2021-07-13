package events

import (
	"errors"
	"fmt"
)

// Executes the Simhospital event generation process
func Execute(config string) (string, error) {
	return fmt.Sprintf("ECHO '%v'", config), errors.New("Unimplemented")
}
