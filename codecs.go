package kupak

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type Object struct {
	data         []byte
	marshaller   func(interface{}) ([]byte, error)
	unmarshaller func([]byte, interface{}) error
}

type Metadata struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
}

func (o *Object) Metadata() (*Metadata, error) {
	return nil
}

func (o *Object) SetLabels(v map[string]string) error {
	return nil
}

func (o *Object) SetAnnotations(v map[string]string) error {
	return nil
}

func (o *Object) AddLabel(kv ...[]string) error {
	return nil
}

func (o *Object) AddAnnotations(kv ...[]string) error {
	return nil
}

func (o *Object) Bytes() ([]byte, error) {
	return nil, nil
}
