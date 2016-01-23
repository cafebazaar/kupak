package kupak

import (
	//	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

type Object struct {
	data         []byte
	marshaller   func(interface{}) ([]byte, error)
	unmarshaller func([]byte, interface{}) error
	kind         string
	value        interface{}
}

type Metadata struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`
}

func NewObject(data []byte,
	marshaller func(interface{}) ([]byte, error),
	unmarshaller func([]byte, interface{}) error) (*Object, error) {
	obj := &Object{
		data:         data,
		marshaller:   marshaller,
		unmarshaller: unmarshaller,
	}
	if err := obj.unmarshallToValue(); err != nil {
		return nil, err
	}
	return obj, nil
}

func (obj *Object) unmarshallToValue() error {
	// what's the kind
	metadata, err := obj.Metadata()
	if err != nil {
		return err
	}
	if metadata.Kind == "" {
		obj.kind = "unknown"
	} else {
		obj.kind = metadata.Kind
	}

	// unmarshall to value if it's a known type
	switch obj.kind {
	case "ReplicationController":
		v := v1.ReplicationController{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Pod":
		v := v1.Pod{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Service":
		v := v1.Service{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "PersistentVolume":
		v := v1.PersistentVolume{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Secret":
		v := v1.Service{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Namespace":
		v := v1.Namespace{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "ServiceAccount":
		v := v1.ServiceAccount{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "DaemonSet":
		v := v1beta1.DaemonSet{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Job":
		v := v1beta1.Job{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Ingress":
		v := v1beta1.Ingress{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	case "Deployment":
		v := v1beta1.Deployment{}
		obj.unmarshaller(obj.data, &v)
		obj.value = &v
	default:
		obj.value = nil
	}
	return nil
}

func (o *Object) Bytes() ([]byte, error) {
	return o.data, nil
}

func (o *Object) Metadata() (*Metadata, error) {
	meta := Metadata{}
	return &meta, o.unmarshaller(o.data, &meta)
}

func (o *Object) SetLabels(labels map[string]string) error {

}

func (o *Object) SetAnnotations(Annotations map[string]string) error {

}

func (o *Object) SetInnerPodLabels(labels map[string]string) error {

}

func (o *Object) SetInnerPodAnnotations(Annotations map[string]string) error {

}

func (o *Object) InnerPodTemplateMetadata() (*Metadata, error) {
	var spec *v1.PodTemplateSpec
	switch o.kind {
	case "ReplicationController":
		spec = o.ReplicationController().Spec.Template
	case "DaemonSet":
		spec = o.DaemonSet().Spec.Template
	case "Job":
		spec = &o.Job().Spec.Template
	case "Deployment":
		spec = o.Deployment().Spec.Template
	default:
		return nil, nil
	}
	return &Metadata{
		ObjectMeta: spec.ObjectMeta,
	}, nil
}

func (o *Object) Kind() string {
	return o.kind
}

func (o *Object) ReplicationController() *v1.ReplicationController {
	return o.value.(*v1.ReplicationController)
}

func (o *Object) Pod() *v1.Pod {
	return o.value.(*v1.Pod)
}

func (o *Object) Service() *v1.Service {
	return o.value.(*v1.Service)
}

func (o *Object) PersistentVolume() *v1.PersistentVolume {
	return o.value.(*v1.PersistentVolume)
}

func (o *Object) Secret() *v1.Secret {
	return o.value.(*v1.Secret)
}

func (o *Object) Namespace() *v1.Namespace {
	return o.value.(*v1.Namespace)
}

func (o *Object) ServiceAccount() *v1.ServiceAccount {
	return o.value.(*v1.ServiceAccount)
}

func (o *Object) DaemonSet() *v1beta1.DaemonSet {
	return o.value.(*v1beta1.DaemonSet)
}

func (o *Object) Job() *v1beta1.Job {
	return o.value.(*v1beta1.Job)
}

func (o *Object) Ingress() *v1beta1.Ingress {
	return o.value.(*v1beta1.Ingress)
}

func (o *Object) Deployment() *v1beta1.Deployment {
	return o.value.(*v1beta1.Deployment)
}
