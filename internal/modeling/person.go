package modeling

import (
	"math"
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
)

type Person struct {
	Identifier []config.Code `yaml:"identifier,omitempty"`
	Name       Name          `yaml:"name,omitempty"`
	Birthdate  string        `yaml:"birthdate,omitempty"`
	DeathDate  string        `yaml:"death_date,omitempty"`
	Gender     string        `yaml:"gender,omitempty"`
	Address    struct {
		Primary   Address `yaml:"primary,omitempty"`
		Temporary Address `yaml:"temporary,omitempty"`
	}
	Details    map[string]string `yaml:"details,omitempty"`
	Conditions []Condition       `yaml:"conditions,omitempty"`
}

type Name struct {
	Given  []string `yaml:"given,omitempty"`
	Family string   `yaml:"family,omitempty"`
}

type Condition struct {
	Name string `yaml:"name"`
	Code struct {
		System string `yaml:"system"`
		Value  string `yaml:"value"`
	}
	OnsetAge int  `yaml:"onset_age"`
	Terminal bool `yaml:"terminal"`
}

func (person *Person) age() int {
	if person.Birthdate == "" {
		return 0
	}
	dob, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		panic(err)
	}
	return int(math.Floor(time.Now().Sub(dob).Hours() / float64(8760)))
}
