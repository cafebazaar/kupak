package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

/* ==== PackageInfo ==== */

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

/* ==== Package ==== */

type Pak struct {
	PakInfo      `yaml:",inline"`
	Properties   []Property `yaml:"properties,omitempty"`
	TemplateUrls []string   `yaml:"templates"`
	Templates    []string   `yaml:""`
}

type Property struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Default interface{} `yaml:"default,omitempty"`
}

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

/* ==== Repo Index ==== */

type Repo struct {
	Url         string     `yaml:""`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description,omitempty"`
	Maintainer  string     `yaml:"maintainer,omitempty"`
	Index       []*PakInfo `yaml:"packages"`
}

func RepoFromBytes(data []byte) (*Repo, error) {
	var repo Repo
	err := yaml.Unmarshal(data, &repo)
	if err != nil {
		return nil, err
	}
	nameMap := make(map[string]bool)
	for i := range repo.Index {
		if _, has := nameMap[repo.Index[i].Name+":"+repo.Index[i].Version]; has {
			return nil, errors.New("Duplicated package")
		}
		nameMap[repo.Index[i].Name+":"+repo.Index[i].Version] = true
		if repo.Index[i].URL == "" {
			return nil, errors.New("Url doesn't exists or is not correct")
		}
		if repo.Index[i].Version == "" {
			return nil, errors.New("Version doesn't exists or is not correct")
		}
	}
	return &repo, nil
}

func RepoFromUrl(url string) (*Repo, error) {
	data, err := fetchUrl(url)
	if err != nil {
		return nil, err
	}
	repo, err := RepoFromBytes(data)
	if err != nil {
		return nil, err
	}
	repo.Url = url
	return repo, nil
}

func (r *Repo) Has(pak string) bool {
	for i := range r.Index {
		if r.Index[i].Name == pak {
			return true
		}
	}
	return false
}

func (r *Repo) HasVersion(pak string, version string) bool {
	for i := range r.Index {
		if r.Index[i].Name == pak && r.Index[i].Version == version {
			return true
		}
	}
	return false
}

func (r *Repo) Pak(pak string, version string) (*Pak, error) {
	for i := range r.Index {
		if r.Index[i].Name == pak && r.Index[i].Version == version {
			url := path.Join(path.Dir(r.Url), r.Index[i].URL)
			fmt.Println(url + "===")
			pak, err := PakFromUrl(url)
			if err != nil {
				return nil, err
			}
			return pak, nil
		}
	}
	return nil, errors.New("Package not found")
}

func fetchUrl(url string) ([]byte, error) {
	if strings.HasPrefix(strings.ToLower(url), "http://") ||
		strings.HasPrefix(strings.ToLower(url), "https://") {
		c := &http.Client{}
		resp, err := c.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return ioutil.ReadFile(url)
	}
}

func main() {
	repo, err := RepoFromUrl("paks/index.yaml")
	if err != nil {
		panic(err)
	}
	for i := range repo.Index {
		fmt.Println(repo.Index[i].String())
	}

	pak, err := repo.Pak("test", "1.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(pak.Name)
}
