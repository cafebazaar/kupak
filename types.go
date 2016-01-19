package kupak

import (
	"strings"
	"text/template"
)

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

func (p *PakInfo) FormatedString() string {
	return strings.Join([]string{"hi"}, "\n")
}

type Pak struct {
	PakInfo      `yaml:",inline"`
	Properties   []Property `yaml:"properties,omitempty"`
	ResourceUrls []string   `yaml:"resources"`

	// Populated from resources
	Templates []*template.Template `yaml:""`
}

type Property struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Description string      `yaml:"description,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}

type Repo struct {
	Url         string     `yaml:""`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description,omitempty"`
	Maintainer  string     `yaml:"maintainer,omitempty"`
	Index       []*PakInfo `yaml:"packages"`
}
