package modeling

import (
	"math"
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
)

type Diagnosis struct {
	Name            string   `yaml:"name"`
	Code            Code     `yaml:"code"`
	Commonality     float32  `yaml:"commonality"`
	OnsetAges       []string `yaml:"onset_ages,omitempty"`
	Mortality       float32  `yaml:"mortality"`
	MortalityWindow string   `yaml:"mortality_window"`
	Nature          string   `yaml:"nature"`
}

func (diag *Diagnosis) Manifest(patients []Patient) int {
	ages := make([][2]int, len(diag.OnsetAges))
	for i, r := range diag.OnsetAges {
		low, high := config.SplitRange(r)
		if high < 0 {
			high = config.MAX_AGE
		}
		ages[i] = [2]int{low, high}
	}
	pats := findPats(diag.Commonality, ages, patients)
	for p := 0; p < len(pats); p++ {
		if patients[p].Conditions == nil {
			patients[p].Conditions = make([]Condition, 0)
		}
		condition := Condition{}
		condition.Name = diag.Name
		condition.Code = diag.Code

		onsetIndex := config.RandInt(0, len(diag.OnsetAges)-1)
		l, h := config.SplitRange(diag.OnsetAges[onsetIndex])
		if h == -1 {
			h = config.MAX_AGE
		}
		patAge := patients[p].Demographics.Age()
		if h > patAge {
			h = patAge
		}
		condition.OnsetAge = config.RandInt(l, h)

		condition.Terminal = config.RandFloat() <= diag.Mortality
		if condition.Terminal {
			window := config.WindowRange(diag.MortalityWindow)
			start, err := time.Parse("2006-01-02", patients[p].Demographics.Birthdate)
			start = start.AddDate(condition.OnsetAge, 0, config.RandInt(0, 365))
			if err != nil {
				panic(err)
			}
			death := start.Add(window)
			patients[p].Demographics.DeathDate = death.Format("2006-01-02")
		}
		patients[p].Conditions = append(patients[p].Conditions, condition)
	}
	return len(pats)
}

func findPats(commonality float32, ages [][2]int, patients []Patient) []int {
	popSize := len(patients)
	expected := float64(commonality * float32(popSize))
	cases := int(math.Floor(expected))
	if cases < 1 {
		roll := config.RandFloat64()
		if roll < expected {
			cases = 1
		}
	}

	possiblePats := make([]int, 0)
	for i := 0; i < popSize; i++ {
		if patients[i].Demographics.DeathDate != "" {
			continue
		}
		dob, err := time.Parse("2006-01-02", patients[i].Demographics.Birthdate)
		if err != nil {
			panic(err)
		}
		age := int(math.Floor(time.Now().Sub(dob).Hours() / float64(8760)))
		for _, r := range ages {
			if r[0] <= age && age < r[1] {
				possiblePats = append(possiblePats, i)
				break
			}
		}
	}

	pats := make([]int, cases)
	for i := 0; i < cases; i++ {
		l := len(possiblePats)
		roll := config.RandInt(0, l)
		pats[i] = possiblePats[roll]

		possiblePats[roll] = possiblePats[l-1]
		possiblePats = possiblePats[:len(possiblePats)-1]
	}

	return pats
}
