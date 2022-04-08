package modeling

type Condition struct {
	Name      string            `yaml:"name"`
	Code      Code              `yaml:"code"`
	OnsetAge  int               `yaml:"onset_age"`
	OnsetDate string            `yaml:"onset_date"`
	Terminal  bool              `yaml:"terminal,omitempty"`
	Type      string            `yaml:"type"`
	Details   map[string]string `yaml:"details"`
}
