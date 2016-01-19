package kupak

import (
	"errors"
	"gopkg.in/yaml.v2"
)

func validateProperties(properties []Property) error {
	nameMap := make(map[string]bool)
	for i := range properties {
		if _, has := nameMap[properties[i].Name]; has {
			return errors.New("Duplicated property")
		}
		// validating types
		switch properties[i].Type {
		case "int":
		case "number":
		case "string":
			// TODO validate the default value and other type specification
			_ = "ok"
		default:
			return errors.New("Specified type is not valid")
		}
	}
	return nil
}

func PakFromUrl(url string) (*Pak, error) {
	data, err := fetchUrl(url)
	if err != nil {
		return nil, err
	}
	pak := Pak{}
	if err := yaml.Unmarshal(data, &pak); err != nil {
		return nil, err
	}
	if err := validateProperties(pak.Properties); err != nil {
		return nil, err
	}
	return &pak, nil
}
