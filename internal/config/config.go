package config

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_AGE     = 110
	UPPER_INDEX = 65
	LOWER_INDEX = 97
)

var RaceMap = make(map[float32]string)
var RaceMapKeys = make([]float32, 0)

var EthnicityMap = make(map[float32]string)
var EthnicityMapKeys = make([]float32, 0)

var SexualityMap = make(map[float32]string)
var SexualityMapKeys = make([]float32, 0)

var GenderIdentityMap = make(map[float32]string)
var GenderIdentityMapKeys = make([]float32, 0)

type Configuration struct {
	Inherit []string `yaml:"inherit"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat() float32 {
	return rand.Float32()
}

func RandInt(l int, h int) int {
	return rand.Intn(h-l) + l
}

func RandMaskValue(mask string) string {
	value := ""
	for _, c := range mask {
		if c == '#' {
			value += fmt.Sprintf("%v", RandInt(0, 10))
		} else if c == '?' {
			switch RandInt(0, 1) {
			case 0:
				value += fmt.Sprintf("%c", (UPPER_INDEX + RandInt(0, 25)))
			case 1:
				value += fmt.Sprintf("%c", (LOWER_INDEX + RandInt(0, 25)))
			}
		} else if c == '*' {
			switch RandInt(0, 2) {
			case 0:
				value += fmt.Sprintf("%v", RandInt(0, 10))
			case 1:
				value += fmt.Sprintf("%c", (UPPER_INDEX + RandInt(0, 25)))
			case 2:
				value += fmt.Sprintf("%c", (LOWER_INDEX + RandInt(0, 25)))
			}

			value += fmt.Sprintf("%v", RandInt(0, 10))
		} else {
			value += string(c)
		}
	}
	return value
}

// General helpers
func SplitRange(rnge string) (int, int) {
	if rnge[len(rnge)-1] == '+' {
		low, err := strconv.ParseInt(strings.TrimRight(rnge, "+"), 0, 64)
		if err != nil {
			panic(fmt.Sprintf("failed to parse low end of range for: %v", err))
		}
		return int(low), -1
	}

	parts := strings.Split(rnge, "-")
	if len(parts) != 2 {
		panic(fmt.Sprintf("malformed range '%v'", rnge))
	}

	low, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 0, 64)
	if err != nil {
		panic(fmt.Sprintf("failed to parse low end of range for: %v", err))
	}

	high, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 0, 64)
	if err != nil {
		panic(fmt.Sprintf("failed to parse high end of range for: %v", err))
	}

	return int(low), int(high)
}
