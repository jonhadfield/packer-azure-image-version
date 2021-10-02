package packer_azure_image_version

import "strings"

// allDefined returns true if all provided strings are not empty
func allDefined(i ...string) bool {
	for x := range i {
		if i[x] == "" {
			return false
		}
	}

	return true
}

// ParseImageDefinitionID accepts an azure resource ID as a string and returns a struct instance containing the components
func ParseImageDefinitionID(rawID string) ImageDefinitionResourceID {
	components := strings.Split(rawID, "/")
	if len(components) != 11 {
		return ImageDefinitionResourceID{}
	}

	return ImageDefinitionResourceID{
		SubscriptionID: components[2],
		ResourceGroup:  components[4],
		Provider:       components[6],
		Gallery:        components[8],
		ImageName:      components[10],
		Raw:            rawID,
	}
}
