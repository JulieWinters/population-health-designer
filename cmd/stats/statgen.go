package stats

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

func Execute(configFile string) (string, error) {

	popStat := ParsePopStatus(configFile)

	patients := make([]modeling.Person, popStat.Rules.Counts.Patients)
	for i := 0; i < popStat.Rules.Counts.Patients; i++ {
		patients[i] = popStat.NewPatient()
	}
	config.Write(patients, popStat.Rules.Output)

	return fmt.Sprintf("ECHO '%v'", configFile), nil
}
