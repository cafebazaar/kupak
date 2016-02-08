package kupak

import (
	"errors"
	"time"

	"github.com/ghodss/yaml"
)

type Object struct {
	data []byte
}

type MetadataMD struct {
	Name                       string            `json:"name,omitempty"`
	GenerateName               string            `json:"generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	SelfLink                   string            `json:"selfLink,omitempty"`
	UID                        string            `json:"uid,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	Generation                 int64             `json:"generation,omitempty"`
	CreationTimestamp          time.Time         `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          *time.Time        `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds *int64            `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty"`
}

type Metadata struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	MetadataMD `json:"metadata,omitempty"`
}

type templateMetadata struct {
	Spec struct {
		Template Metadata `json:"template"`
	} `json:"spec"`
}

func NewObject(data []byte) (*Object, error) {
	obj := &Object{
		data: data,
	}
	return obj, nil
}

func (o *Object) Bytes() ([]byte, error) {
	return o.data, nil
}

func (o *Object) Metadata() (*Metadata, error) {
	meta := Metadata{}
	return &meta, yaml.Unmarshal(o.data, &meta)
}

func (o *Object) SetLabels(labels map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := getMapChild([]string{"metadata"}, m)
	if err != nil {
		return err
	}
	metadata, ok := v.(map[string]interface{})
	if !ok {
		return errors.New("there is no metadata")
	}
	metadata["labels"] = labels
	data, err := yaml.Marshal(m)
	if err == nil {
		o.data = data
	}
	return err
}

func (o *Object) SetAnnotations(annotations map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := getMapChild([]string{"metadata"}, m)
	if err != nil {
		return err
	}
	metadata, ok := v.(map[string]interface{})
	if !ok {
		return errors.New("there is no metadata")
	}
	metadata["annotations"] = annotations
	data, err := yaml.Marshal(m)
	if err == nil {
		o.data = data
	}
	return err
}

func (o *Object) TemplateMetadata() (*Metadata, error) {
	meta := templateMetadata{}
	if err := yaml.Unmarshal(o.data, &meta); err != nil {
		return nil, err
	}
	return &meta.Spec.Template, nil
}

func (o *Object) SetTemplateLabels(labels map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := getMapChild([]string{"spec", "template", "metadata"}, m)
	if err != nil {
		return err
	}
	metadata, ok := v.(map[string]interface{})
	if !ok {
		return errors.New("there is no metadata")
	}
	metadata["labels"] = labels
	data, err := yaml.Marshal(m)
	if err == nil {
		o.data = data
	}
	return err
}

func (o *Object) SetTemplateAnnotations(Annotations map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := getMapChild([]string{"spec", "template", "metadata"}, m)
	if err != nil {
		return err
	}
	metadata, ok := v.(map[string]interface{})
	if !ok {
		return errors.New("there is no metadata")
	}
	metadata["Annotations"] = Annotations
	data, err := yaml.Marshal(m)
	if err == nil {
		o.data = data
	}
	return err
}
