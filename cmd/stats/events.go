package stats

import (
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
)

func (stats *PopStats) createEvents() {
	year := time.Now().AddDate(-1*config.MAX_AGE, 0, 0).Year()
	currentYear := time.Now().Year()

	for year <= currentYear {

		year += 1
	}
}
