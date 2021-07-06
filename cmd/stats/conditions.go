package stats

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

var (
	conditions = make([]Condition, 0)
)

type Condition struct {
	Name string `yaml:"name"`
	Code struct {
		System string `yaml:"system"`
		Value  string `yaml:"value"`
	}
	Commonality float32  `yaml:"commonality"`
	OnsetAges   []string `yaml:"onset_ages,omitempty"`
	Mortality   float32  `yaml:"mortality"`
	Nature      string   `yaml:"nature"`
}

func (condition *Condition) manifest(population *[]modeling.Person) {
	ages := make([][2]int, len(condition.OnsetAges))
	for i, r := range condition.OnsetAges {
		low, high := config.SplitRange(r)
		if high < 0 {
			high = config.MAX_AGE
		}
		ages[i] = [2]int{low, high}
	}
	pats := findPats(condition.Commonality, ages, population)
	fmt.Printf("Patients with %v --> %v", condition.Name, pats)
}

func findPats(commonality float32, ages [][2]int, population *[]modeling.Person) []int {
	popSize := len(*population)
	expected := float64(commonality) * float64(popSize)
	cases := int(math.Floor(expected))
	if cases < 1 {
		roll := rand.Float64()
		if roll < expected {
			cases = 1
		}
	}

	possiblePats := make([]int, 0)
	for i := 0; i < popSize; i++ {
		dob, err := time.Parse("2006-01-02", (*population)[i].Birthdate)
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
		roll := rand.Int31n(int32(len(possiblePats)))
		pats[i] = possiblePats[roll]

		possiblePats[len(possiblePats)-1], possiblePats[i] = possiblePats[i], possiblePats[len(possiblePats)-1]
		possiblePats = possiblePats[:len(possiblePats)-1]
	}

	return pats
}
