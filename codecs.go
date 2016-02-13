package kupak

import (
	"errors"
	"time"

	"github.com/ghodss/yaml"
)

// Object is a represention of a Kubernetes' object
// by using this it is possible to manipulate labels and annotations of the object
type Object struct {
	data []byte
}

// MetadataMD is inlined part of object Metadata which contains labels and annotations
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

// Metadata of a object
type Metadata struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	MetadataMD `json:"metadata,omitempty"`
}

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

// SetAnnotations will replace annotations of the object
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
	v, err := getMapChild([]string{"spec", "template", "metadata"}, m)
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

// Object Status
// ==========================================

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

// PodStatus represents information about the status of a pod. Status may trail the actual
// state of a system.
type PodStatus struct {
	Phase PodPhase `json:"phase,omitempty"`
	// A human readable message indicating details about why the pod is in this state.
	Message string `json:"message,omitempty"`
	// A brief CamelCase message indicating details about why the pod is in this state. e.g. 'OutOfDisk'
	Reason string `json:"reason,omitempty"`

	HostIP string `json:"hostIP,omitempty"`
	PodIP  string `json:"podIP,omitempty"`

	// Date and time at which the object was acknowledged by the Kubelet.
	// This is before the Kubelet pulled the container image(s) for the pod.
	StartTime *time.Time `json:"startTime,omitempty"`
}

// Status returns status field in kubectl get pod xx -o json
func (o *Object) Status() (*PodStatus, error) {
	ps := struct {
		Status *PodStatus `json:"status"`
	}{}
	return ps.Status, yaml.Unmarshal(o.data, &ps)
}
