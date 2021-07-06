package config

import (
	"math/rand"
	"time"
)

var RaceMap = make(map[float32]string)
var RaceMapKeys = make([]float32, 0)

var EthnicityMap = make(map[float32]string)
var EthnicityMapKeys = make([]float32, 0)

var SexualityMap = make(map[float32]string)
var SexualityMapKeys = make([]float32, 0)

var GenderIdentityMap = make(map[float32]string)
var GenderIdentityMapKeys = make([]float32, 0)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandFloat() float32 {
	return rand.Float32()
}
