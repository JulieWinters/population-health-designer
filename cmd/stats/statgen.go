package stats

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

// Executes the population data generation process
func Execute(configFile string) {

	fmt.Printf("Reading configuration file at %v\n", configFile)
	popStat, _ := Parse(configFile)

	patients := make([]modeling.Person, popStat.Rules.Counts.Patients)
	for i := 0; i < popStat.Rules.Counts.Patients; i++ {
		patients[i] = popStat.NewPatient()
	}
	fmt.Printf("Generated %v unique patients\n", len(patients))

	for _, c := range popStat.Diagnoses {
		fmt.Printf("  Manifested %v instances of %v\n", c.Manifest(patients), c.Name)
	}

	for _, p := range patients {
		createEvents(p)
	}

	config.Write(patients, popStat.Rules.Output)
	fmt.Printf("  Writing population details to %v\n", popStat.Rules.Output)
}
