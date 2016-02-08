package kupak

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

type Manager struct {
}

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) Installed(namespace string) ([]*InstalledPak, error) {
	return nil, nil
}

func (m *Manager) Instances(namespace string, pak *Pak) ([]*InstalledPak, error) {
	return nil, nil
}

func (m *Manager) Status(namespace string, instance string) (*InstalledPak, error) {
	return nil, nil
}

// Install a pak with given name
func (m *Manager) Install(pak *Pak, namespace string, properties map[string]interface{}) error {
	rawObjects, err := pak.ExecuteTemplates(properties)
	if err != nil {
		return err
	}
	group, err := uuid.NewV4()
	if err != nil {
		return err
	}
	labels := map[string]string{
		"kupak-group":   group.String(),
		"kupak-pak-url": pak.URL,
	}
	var objects []*Object
	for i := range rawObjects {
		object, err := NewObject(rawObjects[i])
		if err != nil {
			return err
		}
		md, err := object.Metadata()
		if err != nil {
			return err
		}
		mergedLabels := mergeStringMaps(md.Labels, labels)
		err = object.SetLabels(mergedLabels)
		if err != nil {
			return err
		}
		// TODO validation for replication controller - do not ignore
		if templateMd, err := object.TemplateMetadata(); err == nil {
			mergedLabels := mergeStringMaps(templateMd.Labels, labels)
			err = object.SetTemplateLabels(mergedLabels)
			if err != nil {
				return err
			}
		}
		bytes, _ := object.Bytes()
		fmt.Println(string(bytes))
		fmt.Println("----\n----")
		objects = append(objects, object)
	}
	return nil
}

// DeleteInstance will delete a installed pak
func (m *Manager) DeleteInstance(namespace string, group string) ([]*InstalledPak, error) {
	return nil, nil
}
