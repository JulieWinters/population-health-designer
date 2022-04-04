package modeling

type ControlId struct {
	Seed      int `yaml:"seed"`
	Increment int `yaml:"increment"`
	CurrentId int `yaml:"current_id,omitempty"`
}

type EventGen struct {
	Stats     string    `yaml:"stats"`
	ControlId ControlId `yaml:"control_id"`
	Messages  []string  `yaml:"messages"`
	Segments  []string  `yaml:"segments"`
}
