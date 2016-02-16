package kupak

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
)

// TODO refactor err handling and make distinct error types

var (
	KubePath   string
	KubeConfig string
)

func init() {
	if KubePath == "" {
		KubePath = os.Getenv("KUBECTL_PATH")
	}
	if KubePath == "" {
		KubePath = "kubectl" // default value
	}
	if KubeConfig == "" {
		KubeConfig = os.Getenv("KUBECTL_CONFIG")
	}
}

type kubeList struct {
	Items []interface{} `json:"items"`
}

type Kubectl interface {
	// Get returns Objects with given selector
	Get(namespace string, type_ string, selector string) ([]*Object, error)

	// Create creates a kubernetes objects
	Create(namespace string, o *Object) error
}

// KubectlRunner is a real implementation of Kubectl interface which uses kubectl
// executable
type KubectlRunner struct{}

func NewKubectlRunner() (*KubectlRunner, error) {
	return &KubectlRunner{}, nil
}

func (k *KubectlRunner) Get(namespace string, type_ string, selector string) ([]*Object, error) {
	if type_ == "" {
		type_ = "all"
	}
	args := []string{"-o", "json", "get", type_}
	if KubeConfig != "" {
		args = append(args, "--kubeconfig", KubeConfig)
	}
	if selector != "" {
		args = append(args, "-l", selector)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	cmd := exec.Command(KubePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(output))
	}

	list := kubeList{}
	yaml.Unmarshal(output, &list)
	var objects []*Object
	for i := range list.Items {
		data, err := yaml.Marshal(list.Items[i])
		if err != nil {
			return nil, err
		}
		object, err := NewObject(data)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (k *KubectlRunner) Create(namespace string, o *Object) error {
	args := []string{"create", "-f", "-"}
	if KubeConfig != "" {
		args = append(args, "--kubeconfig", KubeConfig)
	}
	if namespace != "" {
		args = append([]string{"--namespace", namespace}, args...)
	}

	cmd := exec.Command(KubePath, args...)

	inBuffer, err := o.Bytes()
	if err != nil {
		return err
	}

	cmd.Stdin = bytes.NewBuffer(inBuffer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}
