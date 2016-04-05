package kubectl

import (
	"errors"

	"github.com/cafebazaar/kupak/pkg/util"
	"github.com/ghodss/yaml"
)

type templateMetadata struct {
	Spec *struct {
		Template *Metadata `json:"template"`
	} `json:"spec"`
}

// NewObject creates an object from given bytes
func NewObject(data []byte) (*Object, error) {
	obj := &Object{
		data: data,
	}
	return obj, nil
}

// Bytes serializes the object
func (o *Object) Bytes() ([]byte, error) {
	return o.data, nil
}

// Metadata returns Metadata, Name and Kind of the object
func (o *Object) Metadata() (*Metadata, error) {
	meta := Metadata{}
	return &meta, yaml.Unmarshal(o.data, &meta)
}

// SetLabels will replace annotations of the object
func (o *Object) SetLabels(labels map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := util.GetMapChild([]string{"metadata"}, m)
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

// SetAnnotations will replace annotations of the object
func (o *Object) SetAnnotations(annotations map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := util.GetMapChild([]string{"metadata"}, m)
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

// TemplateMetadata returns Metadata of objects that contain template like
// ReplicationControl
func (o *Object) TemplateMetadata() (*Metadata, error) {
	meta := templateMetadata{}
	if err := yaml.Unmarshal(o.data, &meta); err != nil {
		return nil, err
	}
	if meta.Spec == nil || meta.Spec.Template == nil {
		return nil, errors.New("template metadata not found")
	}
	return meta.Spec.Template, nil
}

// SetTemplateLabels try to find template labels and replace it with labels
// return an error if there is no template
func (o *Object) SetTemplateLabels(labels map[string]string) error {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(o.data, &m); err != nil {
		return err
	}
	v, err := util.GetMapChild([]string{"spec", "template", "metadata"}, m)
	if err != nil {
		return err
	}
	// TODO if and only if metadata doesn't exists make one (spec and template should exists)
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

// Status returns status field in kubectl get pod xx -o json
func (o *Object) Status() (*PodStatus, error) {
	ps := struct {
		Status *PodStatus `json:"status"`
	}{}
	return ps.Status, yaml.Unmarshal(o.data, &ps)
}
