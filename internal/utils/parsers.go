package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseCardinality(cardinality string) (int, int, error) {
	exp := regexp.MustCompile(`^[0-9]+\.\.([0-9]+|\*)$`)
	if !exp.Match([]byte(cardinality)) {
		return 0, 0, fmt.Errorf("invalid cardinality '%s'", cardinality)
	}
	dot := strings.IndexRune(cardinality, '.')
	min, err := strconv.Atoi(cardinality[:dot])
	if err != nil {
		return 0, 0, err
	}

	var max int
	if strings.Contains(cardinality, "*") {
		max = math.MaxInt32
	} else {
		max, err = strconv.Atoi(cardinality[dot+2:])
		if err != nil {
			return 0, 0, err
		}
	}
	return min, max, nil
}
