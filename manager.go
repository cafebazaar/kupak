package kupak

import (
	"fmt"
	"strings"

	"github.com/nu7hatch/gouuid"
)

type Manager struct {
	kubectl Kubectl
}

func NewManager(kubectl Kubectl) (*Manager, error) {
	return &Manager{kubectl: kubectl}, nil
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
	pakID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	labels := map[string]string{
		"kupak-pak-id": pakID.String(),
		// TODO pak url should be full address with .
		"kupak-pak-url": strings.Replace(pak.URL, "/", "-", -1),
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
		objects = append(objects, object)
	}

	// install
	for i := range objects {
		data, _ := objects[i].Bytes()
		fmt.Println(string(data))
		err := m.kubectl.Create(namespace, objects[i])
		if err != nil {
			// TODO XXXXXXXX rollback
			return err
		}
		fmt.Println("-----")
	}
	return nil
}

// DeleteInstance will delete a installed pak
func (m *Manager) DeleteInstance(namespace string, group string) ([]*InstalledPak, error) {
	return nil, nil
}
