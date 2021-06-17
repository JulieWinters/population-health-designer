package config

import (
	"math/rand"
	"time"
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

func randFloat() float32 {
	return rand.Float32()
}
