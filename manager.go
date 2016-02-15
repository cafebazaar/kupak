package kupak

import (
	"fmt"
	"strings"
	"time"

	"github.com/nu7hatch/gouuid"
)

// Manager manages installation and deploying pak to a kubernetes cluster
type Manager struct {
	kubectl Kubectl
}

// NewManager returns a Manager
func NewManager(kubectl Kubectl) (*Manager, error) {
	return &Manager{kubectl: kubectl}, nil
}

// Installed returns all installed paks in given namespace
func (m *Manager) Installed(namespace string) ([]*InstalledPak, error) {
	objects, err := m.kubectl.Get(namespace, "all", "")
	if err != nil {
		return nil, err
	}

	// group all paks
	groups := make(map[string][]*Object)
	for i := range objects {
		md, err := objects[i].Metadata()
		if err != nil {
			return nil, err
		}
		groupID, has := md.Labels["kp-group-id"]
		if has {
			groups[groupID] = append(groups[groupID], objects[i])
		}
	}

	var installedPaks []*InstalledPak
	for k, v := range groups {
		// create InstalledPak objects from group
		installedPak := &InstalledPak{}
		installedPak.Statuses = make(map[string]*PodStatus)
		installedPak.GroupID = k
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
		}
		installedPaks = append(installedPaks, installedPak)
	}
	return installedPaks, nil
}

// TODO
// func (m *Manager) Instances(namespace string, pak *Pak) ([]*InstalledPak, error)
// func (m *Manager) Status(namespace string, instance string) (*InstalledPak, error)

// Install a pak with given name and returns its groupID
func (m *Manager) Install(pak *Pak, namespace string, properties map[string]interface{}) (string, error) {
	rawObjects, err := pak.ExecuteTemplates(properties)
	if err != nil {
		return "", err
	}

	groupID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	labels := map[string]string{
		"kp-group-id": groupID.String(),
		"kp-pak-id":   pak.ID(),
	}
	annotations := map[string]string{
		"kp-pak-url":      pak.URL,
		"kp-created-time": time.Now().String(),
	}

	var objects []*Object
	for i := range rawObjects {
		object, err := NewObject(rawObjects[i])
		if err != nil {
			return "", err
		}

		md, err := object.Metadata()
		if err != nil {
			return "", err
		}

		mergedLabels := mergeStringMaps(md.Labels, labels)
		if err = object.SetLabels(mergedLabels); err != nil {
			return "", err
		}
		if err = object.SetAnnotations(annotations); err != nil {
			return "", err
		}

		// TODO validation for replication controller - do not ignore
		if templateMd, err := object.TemplateMetadata(); err == nil {
			mergedLabels := mergeStringMaps(templateMd.Labels, labels)
			if err = object.SetTemplateLabels(mergedLabels); err != nil {
				return "", err
			}
		}
		objects = append(objects, object)
	}

	// install
	for i := range objects {
		data, _ := objects[i].Bytes()
		fmt.Println(string(data))
		err := m.kubectl.Create(namespace, objects[i])
		if err != nil {
			// TODO XXXXXXXX rollback
			return groupID.String(), fmt.Errorf("failed calling kubectl.Create: %v", err)
		}
		fmt.Println("-----")
	}
	return groupID.String(), nil
}

// DeleteInstance will delete a installed pak
func (m *Manager) DeleteInstance(namespace string, group string) ([]*InstalledPak, error) {
	return nil, nil
}
