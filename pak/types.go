package pak

import (
	"text/template"

	"git.cafebazaar.ir/alaee/kupak/kubectl"
)

// Info contains basic information about the pak that doesn't need
// to be fetched
type Info struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	URL         string   `json:"url"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Icon        string   `json:"icon,omitempty"`
}

func (p *Info) String() string {
	str := "Pak{" + p.Name + ", Ver: " + p.Version + ", Url: " + p.URL + "}"
	return str
}

// Pak contains all the data and information that needed for installing it
type Pak struct {
	Info         `json:",inline"`
	Properties   []Property `json:"properties,omitempty"`
	ResourceURLs []string   `json:"resources"`

	// Populated from resources
	Templates []*template.Template `json:"-"`
}

// Property contains definition of every property that required for generating
// pak templates
type Property struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// InstalledPak Represents an installed pak with a unique Group
type InstalledPak struct {
	Group            string
	Namespace        string
	PakURL           string
	PropertiesValues map[string]interface{}
	Objects          []*kubectl.Object

	// Map of pod's name and its status
	Statuses map[string]*kubectl.PodStatus
}

// Repo represents an index file that contains list of paks
type Repo struct {
	URL         string  `json:""`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Maintainer  string  `json:"maintainer,omitempty"`
	Paks        []*Info `json:"packages"`
}
