package manager

import (
	"fmt"

	"git.cafebazaar.ir/alaee/kupak/pkg/kubectl"
)

// Manager manages installation and deploying pak to a kubernetes cluster
type Manager struct {
	kubectl kubectl.Kubectl
}

// NewManager returns a Manager
func NewManager(kubectl kubectl.Kubectl) (*Manager, error) {
	return &Manager{kubectl: kubectl}, nil
}

// HasGroup checks is the specfied group is unique or not
func (m *Manager) HasGroup(namespace string, group string) (bool, error) {
	objects, err := m.kubectl.Get(namespace, "all", "kp-group="+group)
	if err != nil {
		return true, fmt.Errorf("HasGroup: %v", err)
	}
	if len(objects) > 0 {
		return true, nil
	}
	return false, nil
}
