package kupak

import (
	"errors"

	"gopkg.in/yaml.v2"
)

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

func RepoFromURL(url string) (*Repo, error) {
	data, err := fetchURL(url)
	if err != nil {
		return nil, err
	}
	repo, err := RepoFromBytes(data)
	if err != nil {
		return nil, err
	}
	repo.URL = url
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
			url := joinURL(r.URL, r.Index[i].URL)
			pak, err := PakFromURL(url)
			if err != nil {
				return nil, err
			}
			return pak, nil
		}
	}
	return nil, errors.New("Package not found")
}
