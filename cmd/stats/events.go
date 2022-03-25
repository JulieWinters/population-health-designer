package stats

import (
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

func createEvents(patient modeling.Person) {
	year := time.Now().AddDate(-1*config.MAX_AGE, 0, 0).Year()
	currentYear := time.Now().Year()

	for year <= currentYear {
		//conditions := patient.ConditionsBetween(fmt.Sprintf("%v-01-01", year), fmt.Sprintf("%v-12-31", year))
		year += 1
	}
}
