package kupak

import (
	"errors"

	"github.com/ghodss/yaml"
)

// RepoFromBytes make *Repo from data
func RepoFromBytes(data []byte) (*Repo, error) {
	var repo Repo
	err := yaml.Unmarshal(data, &repo)
	if err != nil {
		return nil, err
	}
	nameMap := make(map[string]bool)
	for i := range repo.Paks {
		if _, has := nameMap[repo.Paks[i].Name+":"+repo.Paks[i].Version]; has {
			return nil, errors.New("Duplicated package")
		}
		nameMap[repo.Paks[i].Name+":"+repo.Paks[i].Version] = true
		if repo.Paks[i].URL == "" {
			return nil, errors.New("Url doesn't exists or is not correct")
		}
		if repo.Paks[i].Version == "" {
			return nil, errors.New("Version doesn't exists or is not correct")
		}
	}
	return &repo, nil
}

// RepoFromURL fetches index file specified by url and returns a *Repo
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

// Has checks if repo contains the pak
func (r *Repo) Has(pak string) bool {
	for i := range r.Paks {
		if r.Paks[i].Name == pak {
			return true
		}
	}
	return false
}

// HasVersion checks if repo contains the pak with specific version
func (r *Repo) HasVersion(pak string, version string) bool {
	for i := range r.Paks {
		if r.Paks[i].Name == pak && r.Paks[i].Version == version {
			return true
		}
	}
	return false
}

// Pak finds a pak with specified version and returns it
func (r *Repo) Pak(pak string, version string) (*Pak, error) {
	for i := range r.Paks {
		if r.Paks[i].Name == pak && r.Paks[i].Version == version {
			url := joinURL(r.URL, r.Paks[i].URL)
			pak, err := PakFromURL(url)
			if err != nil {
				return nil, err
			}
			return pak, nil
		}
	}
	return nil, errors.New("Package not found")
}
