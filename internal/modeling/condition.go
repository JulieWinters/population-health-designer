package modeling

type Condition struct {
	Name     string `yaml:"name"`
	Code     Code   `yaml:"code"`
	OnsetAge int    `yaml:"onset_age"`
	Terminal bool   `yaml:"terminal"`
	Type     string `yaml:"type"`
}
