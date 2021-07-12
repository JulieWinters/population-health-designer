module github.com/JulieWinters/population-health-designer

go 1.14

require (
	github.com/go-delve/delve v1.6.1
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/JulieWinters/population-health-designer/cmd/stats => ./cmd/stats

replace github.com/JulieWinters/population-health-designer/cmd/events => ./cmd/events

replace github.com/JulieWinters/population-health-designer/cmd/maps => ./cmd/maps

replace github.com/JulieWinters/population-health-designer/internal/config => ./internal/config

replace github.com/JulieWinters/population-health-designer/internal/modeling => ./internal/modeling
