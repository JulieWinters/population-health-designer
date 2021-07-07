package modeling

type Person struct {
	Identifier []Code `yaml:"identifier,omitempty"`
	Name       Name   `yaml:"name,omitempty"`
	Birthdate  string `yaml:"birthdate,omitempty"`
	Address    struct {
		Primary   Address `yaml:"primary,omitempty"`
		Temporary Address `yaml:"temporary,omitempty"`
	}
	Details map[string]string `yaml:"details,omitempty"`
}

type Name struct {
	Given  []string `yaml:"given,omitempty"`
	Family string   `yaml:"family,omitempty"`
}

type Code struct {
	System string `yaml:"system,omitempty"`
	Value  string `yaml:"value,omitempty"`
}
