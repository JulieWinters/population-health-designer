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

func (person *Person) Age() int {
	return person.AgeAt(time.Now().Format("2006-01-02"))
}

func (person *Person) AgeAt(date string) int {
	if person.Birthdate == "" {
		return 0
	}
	dob, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		panic(err)
	}
	when, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return int(math.Floor(when.Sub(dob).Hours() / float64(8760)))
}

func (person *Person) ConditionsBetween(start string, end string) []*Condition {

	found := make([]*Condition, 0)

	ageStart := person.AgeAt(start)
	ageEnd := person.AgeAt(end)

	for _, cond := range person.Conditions {
		if ageEnd <= cond.OnsetAge || cond.OnsetAge <= ageStart {
			found = append(found, &cond)
		}
	}

	return found
}
