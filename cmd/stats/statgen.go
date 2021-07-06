package stats

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

func Execute(configFile string) (string, error) {

	// var popStat PopStats
	// popStat.Parse(configFile)
	popStat := Parse(configFile)

	// fmt.Printf("Random race for new patient: %v\n", popStat.Distributions.RandRace())
	// fmt.Printf("Random ethnicity for new patient: %v\n", popStat.Distributions.RandEthnicity())
	// fmt.Printf("Random sexuality for new patient: %v\n", popStat.Distributions.RandSexuality())
	// fmt.Printf("Random gender identity for new patient: %v\n", popStat.Distributions.RandGender())
	// fmt.Printf("Random name for new patient: %v\n", popStat.Names.RandFeminine())
	// fmt.Printf("Random address for new patient: %v\n", popStat.Addresses.RandResidential())

	patients := make([]modeling.Person, popStat.Rules.Counts.Patients)
	for i := 0; i < popStat.Rules.Counts.Patients; i++ {
		patients[i] = popStat.NewPatient()
	}

	config.Write(patients, popStat.Rules.Output)

	return fmt.Sprintf("ECHO '%v'", configFile), nil
}
