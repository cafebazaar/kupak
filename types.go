package kupak

import "text/template"

// PakInfo contains basic information about the pak that doesn't need
// to be fetched
type PakInfo struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

func (p *PakInfo) String() string {
	str := "Pak{" + p.Name + ", Ver: " + p.Version + ", Url: " + p.URL + "}"
	return str
}

// Pak contains all the data and information that needed for installing it
type Pak struct {
	PakInfo      `yaml:",inline"`
	Properties   []Property `yaml:"properties,omitempty"`
	ResourceUrls []string   `yaml:"resources"`

	// Populated from resources
	Templates []*template.Template `yaml:""`
}

// Property contains definition of every property that required for generating
// pak templates
type Property struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Description string      `yaml:"description,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}

// Repo represents an index file that contains list of paks
type Repo struct {
	URL         string     `yaml:""`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description,omitempty"`
	Maintainer  string     `yaml:"maintainer,omitempty"`
	Index       []*PakInfo `yaml:"packages"`
}

// Status represents current state of an installed pak
type Status int

const (
	StatusError Status = iota
	StatusRunning
	StatusDeleting
)

// InstalledPak Represents an installed pak with a unique Group
type InstalledPak struct {
	Group            string
	Namespace        string
	PakURL           string
	PropertiesValues map[string]interface{}
	Objects          []interface{}
	Status           Status
}
