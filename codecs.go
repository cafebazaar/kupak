package kupak

import (
	//	"k8s.io/kubernetes/pkg/api"
	"errors"

	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

type Object struct {
	data []byte
}

type Metadata struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`
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
	v, err := GetMapChild([]string{"metadata"}, m)
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
	v, err := GetMapChild([]string{"metadata"}, m)
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
	v, err := GetMapChild([]string{"spec", "template", "metadata"}, m)
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
	v, err := GetMapChild([]string{"spec", "template", "metadata"}, m)
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
