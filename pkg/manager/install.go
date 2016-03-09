package manager

import (
	"fmt"
	"time"

	"git.cafebazaar.ir/alaee/kupak/pkg/kubectl"
	"git.cafebazaar.ir/alaee/kupak/pkg/pak"
	"git.cafebazaar.ir/alaee/kupak/pkg/util"
	"gopkg.in/yaml.v2"
)

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
