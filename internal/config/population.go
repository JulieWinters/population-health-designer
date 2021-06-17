package config

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

type AddressTemplate struct {
	City                string  `yaml:"city"`
	State               string  `yaml:"state"`
	PostalMask          string  `yaml:"postal_mask"`
	StreetRange         string  `yaml:"street_range"`
	BuildingNumberRange string  `yaml:"building_number_range"`
	AptPercentage       float32 `yaml:"apt_percentage"`
}

type Addresses struct {
	Residential []AddressTemplate
	Commercial  []AddressTemplate
}

type Names struct {
	Masculine []modeling.Name `yaml:"masculine"`
	Feminine  []modeling.Name `yaml:"feminine"`
}

type Distributions struct {
	Race           map[string]float32 `yaml:"race"`
	Ethnicity      map[string]float32 `yaml:"ethnicity"`
	Sexuality      map[string]float32 `yaml:"sexuality"`
	GenderIdentity map[string]float32 `yaml:"gender_identity"`
}

type Rules struct {
	Counts struct {
		Patients  int `yaml:"patients"`
		Providers int `yaml:"providers"`
	}
	Output string `yaml:"output"`
}

type PopStats struct {
	Inherit       []string `yaml:"inherit"`
	Rules         Rules    `yaml:"rules"`
	Distributions Distributions
	Names         Names
	Addresses     Addresses
}

func (stats *PopStats) NewPatient() modeling.Person {

	patient := modeling.Person{}
	patient.Details = make(map[string]string)

	patient.Name = stats.Names.RandMasculine()
	patient.Address.Primary = stats.Addresses.RandResidential()

	temp := stats.Distributions.RandGender()
	if temp != "" {
		patient.Details["gender"] = temp
	}

	temp = stats.Distributions.RandSexuality()
	if temp != "" {
		patient.Details["sexuality"] = temp
	}

	temp = stats.Distributions.RandEthnicity()
	if temp != "" {
		patient.Details["ethnicity"] = temp
	}

	temp = stats.Distributions.RandRace()
	if temp != "" {
		patient.Details["race"] = temp
	}

	return patient
}

// functions for generating demographic data
func (dist *Distributions) RandRace() string {
	return randDistItem(dist.Race, raceMap, &raceMapKeys, "")
}

func (dist *Distributions) RandEthnicity() string {
	return randDistItem(dist.Ethnicity, ethnicityMap, &ethnicityMapKeys, "")
}

func (dist *Distributions) RandSexuality() string {
	return randDistItem(dist.Sexuality, sexualityMap, &sexualityMapKeys, "heterosexual")
}

func (dist *Distributions) RandGender() string {
	return randDistItem(dist.GenderIdentity, genderIdentityMap, &genderIdentityMapKeys, "")
}

func randDistItem(source map[string]float32, distMap map[float32]string, distKeys *[]float32, defaultValue string) string {
	if len(distMap) == 0 {
		*distKeys = make([]float32, len(source))
		last := float32(0)
		count := 0
		for k, v := range source {
			last += v
			distMap[last] = k
			(*distKeys)[count] = last
			count++
		}
		sort.Slice(*distKeys, func(i, j int) bool { return (*distKeys)[i] < (*distKeys)[j] })
		if last > float32(1) {
			panic("a configured distribution was not <= 1.0")
		}
	}

	f := randFloat()
	lowerBound := float32(0)
	upperBound := lowerBound
	for i := 0; i < len(*distKeys); i++ {
		upperBound = (*distKeys)[i]
		if lowerBound <= f && f < upperBound {
			return distMap[(*distKeys)[i]]
		}
		lowerBound = upperBound
	}
	return defaultValue
}

// functions for generating names
func (name *Names) RandMasculine() modeling.Name {
	return name.Masculine[rand.Intn(len(name.Masculine))]
}

func (name *Names) RandFeminine() modeling.Name {
	l := len(name.Feminine)
	r := rand.Intn(l)
	return name.Feminine[r]
}

// functions for generating addresses
func (addr *Addresses) RandResidential() modeling.Address {
	return randAddress(addr.Residential)
}

func (addr *Addresses) RandCommercial() modeling.Address {
	return randAddress(addr.Commercial)
}

func randAddress(addrs []AddressTemplate) modeling.Address {
	tmp := addrs[rand.Intn(len(addrs))]

	randAddr := modeling.Address{}
	randAddr.City = tmp.City
	randAddr.State = tmp.State

	for _, c := range tmp.PostalMask {
		if c == '*' {
			randAddr.PostalCode += fmt.Sprintf("%v", rand.Intn(10))
		} else {
			randAddr.PostalCode += string(c)
		}
	}

	low, high := splitRange(tmp.BuildingNumberRange)
	building := rand.Intn(high-low) + low
	randAddr.BuildingNumber = fmt.Sprint(building)

	low, high = splitRange(tmp.StreetRange)
	street := rand.Intn(high-low) + low
	var streetRank string
	switch street {
	case 1:
		streetRank = "st"
	case 2:
		streetRank = "nd"
	case 3:
		streetRank = "rd"
	default:
		streetRank = "th"
	}

	var streetType string
	if rand.Intn(2) == 0 {
		streetType = "Street"
	} else {
		streetType = "Avenue"
	}

	randAddr.Street = append(randAddr.Street, fmt.Sprintf("%v %v%v %v", building, street, streetRank, streetType))

	aptRoll := rand.Float32()
	if aptRoll <= tmp.AptPercentage {
		aptNum := rand.Intn(99) + 1
		randAddr.Street = append(randAddr.Street, fmt.Sprintf("Apt. %v", aptNum))
	}

	return randAddr
}

// General helpers
func splitRange(rnge string) (int, int) {
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
