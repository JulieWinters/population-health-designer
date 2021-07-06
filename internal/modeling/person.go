package modeling

type Person struct {
	Identifier string `yaml:"string"`
	Name       Name   `yaml:"name,omitempty"`
	Birthdate  string `yaml:"birthdate,omitempty"`
	Address    struct {
		Primary   Address `yaml:"primary,omitempty"`
		Temporary Address `yaml:"temporary,omitempty"`
	}
	Details map[string]string
}

type Name struct {
	Given  []string `yaml:"given,omitempty"`
	Family string   `yaml:"family,omitempty"`
}
