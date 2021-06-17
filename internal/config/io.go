package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func Write(data interface{}, file string) error {

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Parse a PopStats structure from the yaml
func Parse(file string) (PopStats, error) {
	var out PopStats

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return out, err
	}

	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return out, err
	}

	if len(out.Inherit) > 0 {
		fmt.Printf("Inheritance for %v: %v\n", file, len(out.Inherit))
		for i := len(out.Inherit) - 1; i >= 0; i-- {

			// Parse the inherited config
			parent, err := Parse(out.Inherit[i])
			if err != nil {
				return out, err
			}
			mergeConfigs(&out, &parent)
		}
	}

	return out, nil
}

/*
	This merges the parent fields that are populated into the child iff the child's
	version of the field is empty. This ensures that anything set in the child isn't
	overwritten during the merge
*/
func mergeConfigs(child *PopStats, parent *PopStats) {
	// Distributions
	if parent.Distributions.Ethnicity != nil && child.Distributions.Ethnicity == nil {
		child.Distributions.Ethnicity = parent.Distributions.Ethnicity
	}
	if parent.Distributions.Race != nil && child.Distributions.Race == nil {
		child.Distributions.Race = parent.Distributions.Race
	}
	if parent.Distributions.Sexuality != nil && child.Distributions.Sexuality == nil {
		child.Distributions.Sexuality = parent.Distributions.Sexuality
	}
	if parent.Distributions.GenderIdentity != nil && child.Distributions.GenderIdentity == nil {
		child.Distributions.GenderIdentity = parent.Distributions.GenderIdentity
	}

	// Names
	if len(parent.Names.Masculine) != 0 && len(child.Names.Masculine) == 0 {
		child.Names.Masculine = parent.Names.Masculine
	}

	if len(parent.Names.Feminine) != 0 && len(child.Names.Feminine) == 0 {
		child.Names.Feminine = parent.Names.Feminine
	}

	// Address
	if len(parent.Addresses.Residential) != 0 && len(child.Addresses.Residential) == 0 {
		child.Addresses.Residential = parent.Addresses.Residential
	}
	if len(parent.Addresses.Commercial) != 0 && len(child.Addresses.Commercial) == 0 {
		child.Addresses.Commercial = parent.Addresses.Commercial
	}
}
