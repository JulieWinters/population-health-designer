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

	NANOS_PER_MINUTE = 60000000000
	NANOS_PER_HOUR   = 3600000000000
	NANOS_PER_DAY    = 86400000000000
	NANOS_PER_WEEK   = 604800000000000
	NANOS_PER_MONTH  = 2592000000000000
	NANOS_PER_YEAR   = 31556952000000000
)

var RaceMap = make(map[float32]string)
var RaceMapKeys = make([]float32, 0)

var EthnicityMap = make(map[float32]string)
var EthnicityMapKeys = make([]float32, 0)

var SexualityMap = make(map[float32]string)
var SexualityMapKeys = make([]float32, 0)

var GenderIdentityMap = make(map[float32]string)
var GenderIdentityMapKeys = make([]float32, 0)

var AgeMap = make(map[float32]string)
var AgeMapKeys = make([]float32, 0)

type Code struct {
	System string `yaml:"system,omitempty"`
	Value  string `yaml:"value,omitempty"`
}

type Configuration struct {
	Inherit []string `yaml:"inherit"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat() float32 {
	return rand.Float32()
}

func RandFloat64() float64 {
	return rand.Float64()
}

func RandInt(l int, h int) int {
	if h < l {
		return 0
	} else if h == l {
		return l
	}
	return rand.Intn(h-l) + l
}

func RandInt64(l int64, h int64) int64 {
	if h < l {
		return 0
	} else if h == l {
		return l
	}
	return rand.Int63n(h-l) + l
}

func RandMaskedTime(input string, mask string) string {
	value := strings.Builder{}
	for i := 0; i < len(mask); i++ {
		if mask[i] == '_' {
			if i < len(input) {
				value.WriteByte(input[i])
			}
		} else {
			if i+1 < len(mask) {
				if mask[i] == 'H' && mask[i+1] == 'H' {
					hour := strconv.Itoa(RandInt(0, 24))
					i++
					value.WriteString(hour)
				} else if mask[i] == 'h' && mask[i+1] == 'h' {
					hour := strconv.Itoa(RandInt(0, 12) + 1)
					i++
					value.WriteString(hour)
				} else if mask[i] == 'm' && mask[i+1] == 'm' {
					min := strconv.Itoa(RandInt(0, 59))
					i++
					value.WriteString(min)
				} else if mask[i] == 's' && mask[i+1] == 's' {
					sec := strconv.Itoa(RandInt(0, 59))
					i++
					value.WriteString(sec)
				} else {
					value.WriteByte(mask[i])
				}
			}
		}
	}
	return value.String()
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

func WindowRange(window string) time.Duration {
	var dur time.Duration
	dur = -1
	if window[len(window)-1] != '+' && !strings.Contains(window, "-") {
		window = strings.TrimRight(window, "+")
		unit := window[len(window)-1]
		quant := window[0 : len(window)-1]
		val, err := strconv.ParseInt(quant, 32, 0)
		if err != nil {
			panic(err)
		}
		switch unit {
		case 'h':
			dur = time.Duration(RandInt64(NANOS_PER_MINUTE, val*NANOS_PER_HOUR))
		case 'd':
			dur = time.Duration(RandInt64(NANOS_PER_HOUR, val*NANOS_PER_DAY))
		case 'w':
			dur = time.Duration(RandInt64(NANOS_PER_DAY, val*NANOS_PER_WEEK))
		case 'm':
			dur = time.Duration(RandInt64(NANOS_PER_WEEK, val*NANOS_PER_MONTH))
		case 'y':
			dur = time.Duration(RandInt64(NANOS_PER_MONTH, val*NANOS_PER_YEAR))
		}
	}

	if dur == -1 {
		panic("Malformed window")
	}
	return dur
}
