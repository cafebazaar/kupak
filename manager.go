package kupak

type Status int

const (
	StatusError Status = iota
	StatusRunning
	StatusDeleting
)

type InstalledPak struct {
	Instance   string
	Namespace  string
	PakUrl     string
	Properties map[string]string
	Objects    []interface{}
	Status     Status
}

type Manager struct {
}

func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

func (m *Manager) Installed(namespace string) ([]*InstalledPak, error) {
	return nil, nil
}

func (m *Manager) Instances(namespace string, pak *Pak) ([]*InstalledPak, error) {
	return nil, nil
}

func (m *Manager) Status(namespace string, instance string) (*InstalledPak, error) {
	return nil, nil
}

// Install a pak with given name
func (m *Manager) Install(pak *Pak, namespace string, instance string, properties map[string]string) error {
	return nil
}

// DeleteInstance will delete a installed pak
func (m *Manager) DeleteInstance(namespace string, instance string) ([]*InstalledPak, error) {
	return nil, nil
}
