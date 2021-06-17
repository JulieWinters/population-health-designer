package modeling

type Person struct {
	Name    Name `yaml:"name,omitempty"`
	Address struct {
		Primary   Address `yaml:"primary,omitempty"`
		Temporary Address `yaml:"temporary,omitempty"`
	}
	Details map[string]string
}

type Name struct {
	Given  []string `yaml:"given,omitempty"`
	Family string   `yaml:"family,omitempty"`
}
