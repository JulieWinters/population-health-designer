package modeling

import (
	"sort"
	"time"

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
	Ages           map[string]float32 `yaml:"ages"`
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
	Inherit       []string      `yaml:"inherit"`
	Rules         Rules         `yaml:"rules"`
	Identifiers   []Identifier  `yaml:"identifiers"`
	Distributions Distributions `yaml:"distributions"`
	Names         Names         `yaml:"names"`
	Addresses     Addresses     `yaml:"addresses"`
	Diagnoses     []Diagnosis   `yaml:"diagnoses"`
}

func (stats *PopStats) NewPatient() Person {

	patient := Person{}
	patient.Identifier = make([]config.Code, 0)
	patient.Details = make(map[string]string)

	if config.RandFloat() > .5 {
		patient.Gender = "F"
		patient.Name = stats.Names.RandFeminine()
	} else {
		patient.Gender = "M"
		patient.Name = stats.Names.RandMasculine()
	}
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

	age := stats.Distributions.RandAge()
	dob := time.Now()
	year := -1 * age
	day := config.RandInt(dob.YearDay(), 366)
	month := -1 * (day / 30)
	day = -1 * (day % 30)
	dob = dob.AddDate(year, month, day)
	patient.Birthdate = dob.Format("2006-01-02")

	for _, id := range stats.Identifiers {
		patId := config.Code{Value: config.RandMaskValue(id.Mask), System: id.Type}
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

func (dist *Distributions) RandAge() int {
	rng := randDistItem(dist.Ages, config.AgeMap, &config.AgeMapKeys, "")
	l, h := config.SplitRange(rng)
	if h < 0 {
		h = config.MAX_AGE
	}
	return config.RandInt(l, h)
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
