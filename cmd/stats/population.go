package stats

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

// Parse a PopStats structure from the yaml
func Parse(file string) (modeling.PopStats, error) {
	fmt.Print("Parsing Population Stats Generation")

	var out modeling.PopStats

	err := config.Parse(file, &out)
	if err != nil {
		return out, err
	}

	if len(out.Inherit) > 0 {

		fmt.Printf("Inheriting from (%v):\n", len(out.Inherit))
		for i := len(out.Inherit) - 1; i >= 0; i-- {

			// Parse the inherited config
			fmt.Printf("  Reading inheritance from %v\n", out.Inherit[i])
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
func mergeConfigs(child *modeling.PopStats, parent *modeling.PopStats) {
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

	// Diagnoses
	if len(parent.Diagnoses) != 0 && len(child.Diagnoses) == 0 {
		child.Diagnoses = parent.Diagnoses
	}
}
