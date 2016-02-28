package manager

import (
	"fmt"
	"strings"
	"time"

	"git.cafebazaar.ir/alaee/kupak/kubectl"
	"git.cafebazaar.ir/alaee/kupak/pak"
	"git.cafebazaar.ir/alaee/kupak/util"
	"github.com/ghodss/yaml"
)

// Manager manages installation and deploying pak to a kubernetes cluster
type Manager struct {
	kubectl kubectl.Kubectl
}

// NewManager returns a Manager
func NewManager(kubectl kubectl.Kubectl) (*Manager, error) {
	return &Manager{kubectl: kubectl}, nil
}

// Installed returns all installed paks in given namespace
func (m *Manager) Installed(namespace string) ([]*pak.InstalledPak, error) {
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

// TODO
// func (m *Manager) Instances(namespace string, pak *Pak) ([]*InstalledPak, error)
// func (m *Manager) Status(namespace string, instance string) (*InstalledPak, error)

// Install a pak with given name and returns its group
func (m *Manager) Install(pak *pak.Pak, namespace string, properties map[string]interface{}) (string, error) {
	var group string
	// try to get group from properties, if it doesn't exists create a random
	// group and add it to properties
	if val, has := properties["group"]; has {
		group = val.(string)
	} else {
		group = util.GenerateRandomGroup()
		properties["group"] = group
	}

	// check for group duplication
	hasGroup, err := m.HasGroup(namespace, group)
	if err != nil {
		return "", err
	}
	if hasGroup {
		return "", fmt.Errorf("Install: duplicated group '%s'", group)
	}

	// execute the templates
	rawObjects, err := pak.ExecuteTemplates(properties)
	if err != nil {
		return "", err
	}

	// apply labels and annotations
	labels := map[string]string{
		"kp-group":  group,
		"kp-pak-id": pak.ID(),
	}
	propertiesYaml, err := yaml.Marshal(properties)
	if err != nil {
		return "", err
	}
	annotations := map[string]string{
		"kp-pak-url":        pak.URL,
		"kp-created-time":   time.Now().String(),
		"kp-pak-properties": string(propertiesYaml),
	}
	var objects []*kubectl.Object
	for i := range rawObjects {
		object, err := kubectl.NewObject(rawObjects[i])
		if err != nil {
			return "", err
		}

		md, err := object.Metadata()
		if err != nil {
			return "", err
		}

		mergedLabels := util.MergeStringMaps(md.Labels, labels)
		if err = object.SetLabels(mergedLabels); err != nil {
			return "", err
		}
		if err = object.SetAnnotations(annotations); err != nil {
			return "", err
		}

		// TODO validation for replication controller - do not ignore
		if templateMd, err := object.TemplateMetadata(); err == nil {
			mergedLabels := util.MergeStringMaps(templateMd.Labels, labels)
			if err = object.SetTemplateLabels(mergedLabels); err != nil {
				return "", err
			}
		}
		objects = append(objects, object)
	}

	// install the objects
	for i := range objects {
		data, _ := objects[i].Bytes()
		fmt.Println(string(data))
		err := m.kubectl.Create(namespace, objects[i])
		if err != nil {
			// TODO XXXXXXXX rollback
			return group, fmt.Errorf("failed calling kubectl.Create: %v", err)
		}
		fmt.Println("-----")
	}

	return group, nil
}

// HasGroup checks is the specfied group is unique or not
func (m *Manager) HasGroup(namespace string, group string) (bool, error) {
	objects, err := m.kubectl.Get(namespace, "all", "kp-group="+group)
	if err != nil {
		return true, fmt.Errorf("HasGroup: %v", err)
	}
	if len(objects) > 0 {
		return true, nil
	}
	return false, nil
}

// DeleteInstance will delete a installed pak
func (m *Manager) DeleteInstance(namespace string, group string) ([]*pak.InstalledPak, error) {
	return nil, nil
}
