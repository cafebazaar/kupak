package pak

import (
	"errors"
	"strings"

	"git.cafebazaar.ir/alaee/kupak/pkg/util"

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
func RepoFromURL(repoURL string) (*Repo, error) {
	// check if index.json or index.yaml is specified in url, if it's not add
	// both one by one and check for it existence
	if !strings.HasSuffix(repoURL, ".json") && !strings.HasSuffix(repoURL, ".yaml") {
		// check if .yaml exists
		yamlURL := util.JoinURL(repoURL, "index.yaml")
		repo, err := RepoFromURL(yamlURL)
		if err == nil {
			return repo, nil
		}

		// check if .json
		jsonURL := util.JoinURL(repoURL, "index.json")
		return RepoFromURL(jsonURL)
	}

	data, err := util.FetchURL(repoURL)
	if err != nil {
		return nil, err
	}
	repo, err := RepoFromBytes(data)
	if err != nil {
		return nil, err
	}
	repo.URL = repoURL
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
			url := util.JoinURL(r.URL, r.Paks[i].URL)
			pak, err := FromURL(url)
			if err != nil {
				return nil, err
			}
			return pak, nil
		}
	}
	return nil, errors.New("Package not found")
}
