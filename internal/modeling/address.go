package modeling

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/config"
)

type Address struct {
	City           string   `yaml:"city,omitempty"`
	State          string   `yaml:"state,omitempty"`
	Street         []string `yaml:"street,omitempty"`
	BuildingNumber string   `yaml:"building_number,omitempty"`
	PostalCode     string   `yaml:"postal_code,omitempty"`
}

// functions for generating addresses
func (addr *Addresses) RandResidential() Address {
	return randAddress(addr.Residential)
}

func (addr *Addresses) RandCommercial() Address {
	return randAddress(addr.Commercial)
}

func randAddress(addrs []AddressTemplate) Address {
	tmp := addrs[config.RandInt(0, len(addrs))]

	randAddr := Address{}
	randAddr.City = tmp.City
	randAddr.State = tmp.State
	randAddr.PostalCode = config.RandMaskValue(tmp.PostalMask)

	low, high := config.SplitRange(tmp.BuildingNumberRange)
	building := config.RandInt(low, high)
	randAddr.BuildingNumber = fmt.Sprint(building)

	low, high = config.SplitRange(tmp.StreetRange)
	street := config.RandInt(0, high-low) + low
	var streetRank string
	switch street {
	case 1:
		streetRank = "st"
	case 2:
		streetRank = "nd"
	case 3:
		streetRank = "rd"
	default:
		streetRank = "th"
	}

	var streetType string
	if config.RandInt(0, 2) == 0 {
		streetType = "Street"
	} else {
		streetType = "Avenue"
	}

	randAddr.Street = append(randAddr.Street, fmt.Sprintf("%v %v%v %v", building, street, streetRank, streetType))

	aptRoll := config.RandFloat()
	if aptRoll <= tmp.AptPercentage {
		aptNum := config.RandInt(0, 99) + 1
		randAddr.Street = append(randAddr.Street, fmt.Sprintf("Apt. %v", aptNum))
	}

	return randAddr
}
