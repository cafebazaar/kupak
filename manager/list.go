package manager

import (
	"strings"

	"git.cafebazaar.ir/alaee/kupak/kubectl"
	"git.cafebazaar.ir/alaee/kupak/pak"
	"gopkg.in/yaml.v2"
)

// List returns all installed paks in given namespace
func (m *Manager) List(namespace string) ([]*pak.InstalledPak, error) {
	objects, err := m.kubectl.Get(namespace, "all", "")
	if err != nil {
		return nil, err
	}

	// group all paks
	groups := make(map[string][]*kubectl.Object)
	for i := range objects {
		md, err := objects[i].Metadata()
		if err != nil {
			return nil, err
		}
		group, has := md.Labels["kp-group"]
		if has {
			groups[group] = append(groups[group], objects[i])
		}
	}

	var installedPaks []*pak.InstalledPak
	for k, v := range groups {
		// create InstalledPak objects from group
		installedPak := &pak.InstalledPak{}
		installedPak.Statuses = make(map[string]*kubectl.PodStatus)
		installedPak.Group = k
		installedPak.Objects = v
		for i := range v {
			md, err := v[i].Metadata()
			if err != nil {
				return nil, err
			}

			// find url
			if url, has := md.Annotations["kp-pak-url"]; installedPak.PakURL == "" && has {
				installedPak.PakURL = url
			}

			// find namespace
			if installedPak.Namespace == "" {
				installedPak.Namespace = md.Namespace
			}

			// extracting all pod statuses
			if strings.ToLower(md.Kind) == "pod" {
				status, err := v[i].Status()
				if err != nil {
					return nil, err
				}
				installedPak.Statuses[md.Name] = status
			}

			// properties values
			if propertiesRaw, has := md.Annotations["kp-pak-properties"]; installedPak.PropertiesValues == nil && has {
				var properties map[string]interface{}
				if err := yaml.Unmarshal([]byte(propertiesRaw), &properties); err != nil {
					return nil, err
				}
				installedPak.PropertiesValues = properties
			}
		}
		installedPaks = append(installedPaks, installedPak)
	}
	return installedPaks, nil
}
