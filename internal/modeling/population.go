package modeling

import (
	"sort"

	"github.com/JulieWinters/population-health-designer/internal/config"
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
	Masculine []Name `yaml:"masculine"`
	Feminine  []Name `yaml:"feminine"`
}

type Distributions struct {
	Race           map[string]float32 `yaml:"race"`
	Ethnicity      map[string]float32 `yaml:"ethnicity"`
	Sexuality      map[string]float32 `yaml:"sexuality"`
	GenderIdentity map[string]float32 `yaml:"gender_identity"`
}

type Identifier struct {
	Mask string `yaml:"mask"`
	Type string `yaml:"type"`
}

type Rules struct {
	Counts struct {
		Patients  int `yaml:"patients"`
		Providers int `yaml:"providers"`
	}
	Output string `yaml:"output"`
}

type PopStats struct {
	Inherit       []string     `yaml:"inherit"`
	Rules         Rules        `yaml:"rules"`
	Identifiers   []Identifier `yaml:"identifiers"`
	Distributions Distributions
	Names         Names
	Addresses     Addresses
}

func (stats *PopStats) NewPatient() Person {

	patient := Person{}
	patient.Identifier = make([]Code, 0)
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

	for _, id := range stats.Identifiers {
		patId := Code{Value: config.RandMaskValue(id.Mask), System: id.Type}
		patient.Identifier = append(patient.Identifier, patId)
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

// functions for generating names
func (name *Names) RandMasculine() Name {
	return name.Masculine[config.RandInt(0, len(name.Masculine))]
}

func (name *Names) RandFeminine() Name {
	l := len(name.Feminine)
	r := config.RandInt(0, l)
	return name.Feminine[r]
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
