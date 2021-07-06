package stats

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
	"gopkg.in/yaml.v2"
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
	return randDistItem(dist.Race, config.RaceMap, &config.RaceMapKeys, "")
}

func (dist *Distributions) RandEthnicity() string {
	return randDistItem(dist.Ethnicity, config.EthnicityMap, &config.EthnicityMapKeys, "")
}

func (dist *Distributions) RandSexuality() string {
	return randDistItem(dist.Sexuality, config.SexualityMap, &config.SexualityMapKeys, "heterosexual")
}

func (dist *Distributions) RandGender() string {
	return randDistItem(dist.GenderIdentity, config.GenderIdentityMap, &config.GenderIdentityMapKeys, "")
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

	f := config.RandFloat()
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

// Parse a PopStats structure from the yaml
// func (stats *PopStats) Parse(file string) error {
func Parse(file string) PopStats {
	var stats PopStats

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &stats)
	if err != nil {
		panic(err)
	}

	//config.Parse(file, stats)

	if len(stats.Inherit) > 0 {
		fmt.Printf("Inheritance for %v: %v\n", file, len(stats.Inherit))
		for i := len(stats.Inherit) - 1; i >= 0; i-- {

			// Parse the inherited config
			// var parent PopStats
			//err := parent.Parse(stats.Inherit[i])
			parent := Parse(stats.Inherit[i])
			if err != nil {
				panic(err)
			}
			mergeConfigs(&stats, &parent)
		}
	}

	return stats
}

/*
	This merges the parent fields that are populated into the child iff the child's
	version of the field is empty. This ensures that anything set in the child isn't
	overwritten during the merge
*/
func mergeConfigs(child *PopStats, parent *PopStats) {
	// Distributions
	if parent.Distributions.Ethnicity != nil && child.Distributions.Ethnicity == nil {
		child.Distributions.Ethnicity = parent.Distributions.Ethnicity
	}
	if parent.Distributions.Race != nil && child.Distributions.Race == nil {
		child.Distributions.Race = parent.Distributions.Race
	}
	if parent.Distributions.Sexuality != nil && child.Distributions.Sexuality == nil {
		child.Distributions.Sexuality = parent.Distributions.Sexuality
	}
	if parent.Distributions.GenderIdentity != nil && child.Distributions.GenderIdentity == nil {
		child.Distributions.GenderIdentity = parent.Distributions.GenderIdentity
	}

	// Names
	if len(parent.Names.Masculine) != 0 && len(child.Names.Masculine) == 0 {
		child.Names.Masculine = parent.Names.Masculine
	}

	if len(parent.Names.Feminine) != 0 && len(child.Names.Feminine) == 0 {
		child.Names.Feminine = parent.Names.Feminine
	}

	// Address
	if len(parent.Addresses.Residential) != 0 && len(child.Addresses.Residential) == 0 {
		child.Addresses.Residential = parent.Addresses.Residential
	}
	if len(parent.Addresses.Commercial) != 0 && len(child.Addresses.Commercial) == 0 {
		child.Addresses.Commercial = parent.Addresses.Commercial
	}
}
