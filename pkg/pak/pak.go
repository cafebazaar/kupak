package pak

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"text/template"

	"git.cafebazaar.ir/alaee/kupak/util"

	"github.com/ghodss/yaml"
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
		case "bool":
		case "string":
			// TODO validate the default value and other type specification
			_ = "ok"
		default:
			return errors.New("Specified type is not valid")
		}
	}
	return nil
}

// ID returns unique id for this pak
func (p *Pak) ID() string {
	md5er := md5.New()
	io.WriteString(md5er, p.URL)
	return fmt.Sprintf("%x", md5er.Sum(nil))
}

func (p *Pak) fetchAndMakeTemplates(baseURL string) error {
	p.Templates = make([]*template.Template, len(p.ResourceURLs))
	for i := range p.ResourceURLs {
		url := util.JoinURL(baseURL, p.ResourceURLs[i])
		data, err := util.FetchURL(url)
		if err != nil {
			return err
		}
		t := template.New(p.ResourceURLs[i])
		t.Delims("$(", ")")
		resourceTemplate, err := t.Parse(string(data))
		if err != nil {
			return err
		}
		p.Templates[i] = resourceTemplate
	}
	return nil
}

// ValidateValues validates given values with corresponding properties
// given values should be contain defaults - use MergeValuesWithDefaults before
// passing values to this function
func (p *Pak) ValidateValues(values map[string]interface{}) error {
	// check all required values are given and their values are ok
	for i := range p.Properties {
		v, has := values[p.Properties[i].Name]
		if !has {
			return errors.New("required property '" + p.Properties[i].Name + "' is not specified")
		}

		ok := false
		switch p.Properties[i].Type {
		case "string":
			_, ok = v.(string)
		case "int":
			_, ok = v.(int)
		case "bool":
			_, ok = v.(bool)
		}
		if !ok {
			return fmt.Errorf("value \"%v\" for property \"%s\" is not correct", v, p.Properties[i].Name)
		}
	}
	return nil
}

// ExecuteTemplates generate resources of a pak with given values
func (p *Pak) ExecuteTemplates(values map[string]interface{}) ([][]byte, error) {
	// merge default values
	values, err := p.MergeValuesWithDefaults(values)
	if err != nil {
		return nil, err
	}

	err = p.ValidateValues(values)
	if err != nil {
		return nil, err
	}
	outputs := make([][]byte, len(p.Templates))
	for i := range p.Templates {
		buffer := &bytes.Buffer{}

		if err := p.Templates[i].Execute(buffer, values); err != nil {
			return nil, err
		}
		outputs[i] = buffer.Bytes()
	}
	return outputs, nil
}

// MergeValuesWithDefaults add default values for properties that not exists in values
func (p *Pak) MergeValuesWithDefaults(values map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for i := range p.Properties {
		value, ok := values[p.Properties[i].Name]
		if ok {
			result[p.Properties[i].Name] = value
		} else if value := p.Properties[i].Default; value != nil {
			values[p.Properties[i].Name] = value
		}
	}
	return result, nil
}

// FromURL reads a pak.yaml file and fetches all the resources files
func FromURL(url string) (*Pak, error) {
	data, err := util.FetchURL(url)
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
	if err := pak.fetchAndMakeTemplates(url); err != nil {
		return nil, err
	}
	return &pak, nil
}
