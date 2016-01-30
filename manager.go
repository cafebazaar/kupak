package kupak

type Status int

const (
	StatusError Status = iota
	StatusRunning
	StatusDeleting
)

type RuntimePak struct {
	Instance   string
	Namespace  string
	PakUrl     string
	Properties map[string]string
	Status     Status
}

type Manager struct {
}

func NewManager() (*Manager, error) {
	return &Manager{}
}

func (m *Manager) Installed(namespace string) ([]*RuntimePak, error) {
	return nil, nil
}

func (m *Manager) Instances(namespace string, pak *Pak) ([]*RuntimePak, error) {
	return nil, nil
}

func (m *Manager) Status(namespace string, instance string) (*RuntimePak, error) {
	return nil, nil
}

func (m *Manager) Install(pak *Pak, namespace string, instance string, properties map[string]string) error {
	return nil
}

func (m *Manager) DeleteInstance(namespace string, instance string) ([]*RuntimePak, error) {
	return nil, nil
}

func (m *Manager) DeleteInstances(namespace string, pak *Pak) ([]*RuntimePak, error) {
	return nil, nil
}
