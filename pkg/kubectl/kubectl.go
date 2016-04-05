package kubectl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/ghodss/yaml"
)

// TODO refactor err handling and make distinct error types

var (
	// KubePath is the path to kubectl executable
	KubePath string

	// KubeConfig is the path to kubeconfig
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
	if KubeConfig == "" {
		KubeConfig = path.Join(os.Getenv("HOME"), ".kube", "config")
	}
}

type kubeList struct {
	Items []interface{} `json:"items"`
}

// KubectlRunner is a real implementation of Kubectl interface which uses kubectl
// executable
type KubectlRunner struct{}

// NewKubectlRunner returns a instance of KubectlRunner that uses kubectl
// external command to implement Kubectl interface
func NewKubectlRunner() (*KubectlRunner, error) {
	return &KubectlRunner{}, nil
}

// Get implements Get of Kubectl interface
func (k *KubectlRunner) Get(namespace string, objType string, selector string) ([]*Object, error) {
	if objType == "" {
		objType = "all"
	}
	args := []string{"-o", "json", "get", objType}
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
		return nil, fmt.Errorf("kubectl error: %v - %s", err, string(output))
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

// Create implements Create of Kubectl interface
func (k *KubectlRunner) Create(namespace string, o *Object) error {
	args := []string{"create", "-f", "-"}
	if KubeConfig != "" {
		args = append(args, "--kubeconfig", KubeConfig)
	}
	if namespace != "" {
		args = append([]string{"--namespace", namespace}, args...)
	}

	inBuffer, err := o.Bytes()
	if err != nil {
		return err
	}

	cmd := exec.Command(KubePath, args...)
	cmd.Stdin = bytes.NewBuffer(inBuffer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl error: %v - %s", err, string(output))
	}
	return nil
}

// Annotate annotate an kubernetes object
func (k *KubectlRunner) Annotate(namespace string, objType string, selector string, annotation string) error {
	args := []string{"annotate", "--overwrite", objType, annotation}
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
		return fmt.Errorf("kubectl error: %v - %s", err, string(output))
	}
    return nil
}

// Label label an kubernetes object
func (k *KubectlRunner) Label(namespace string, objType string, selector string, label string) error {
	args := []string{"label", "--overwrite", objType, label}
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
		return fmt.Errorf("kubectl error: %v - %s", err, string(output))
	}
    return nil
}



// Delete implements Delete of Kubectl interface
func (k *KubectlRunner) Delete(namespace string, o *Object) error {
	args := []string{"delete", "-f", "-"}
	if KubeConfig != "" {
		args = append(args, "--kubeconfig", KubeConfig)
	}
	if namespace != "" {
		args = append([]string{"--namespace", namespace}, args...)
	}

	inBuffer, err := o.Bytes()
	if err != nil {
		return err
	}

	cmd := exec.Command(KubePath, args...)
	cmd.Stdin = bytes.NewBuffer(inBuffer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl error: %v - %s", err, string(output))
	}
	return nil
}
