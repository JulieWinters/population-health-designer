package config

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_AGE = 110
)

var raceMap = make(map[float32]string)
var raceMapKeys = make([]float32, 0)

var ethnicityMap = make(map[float32]string)
var ethnicityMapKeys = make([]float32, 0)

var sexualityMap = make(map[float32]string)
var sexualityMapKeys = make([]float32, 0)

var genderIdentityMap = make(map[float32]string)
var genderIdentityMapKeys = make([]float32, 0)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat() float32 {
	return rand.Float32()
}

func RandInt(l int, h int) int {
	return rand.Intn(h-l) + l
}

type Configuration struct {
	Inherit []string `yaml:"inherit"`
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
