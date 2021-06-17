package modeling

type Address struct {
	City           string   `yaml:"city,omitempty"`
	State          string   `yaml:"state,omitempty"`
	Street         []string `yaml:"street,omitempty"`
	BuildingNumber string   `yaml:"building_number,omitempty"`
	PostalCode     string   `yaml:"postal_code,omitempty"`
}
