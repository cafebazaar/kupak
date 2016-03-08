package manager

import "git.cafebazaar.ir/alaee/kupak/pkg/kubectl"

var removeOrder = []string{"Service", "Deployment", "ReplicationController"}

func stringInArrayIndex(array []string, s string) int {
	for i := range array {
		if array[i] == s {
			return i
		}
	}
	return -1
}

// Remove deletes an installed pak
func (m *Manager) Remove(namespace string, group string) error {
	installedPaks, err := m.listByLabels(namespace, "kp-group="+group)
	if err != nil {
		return err
	}
	// last index is for all other objects
	objectsToDeleteByOrder := make([][]*kubectl.Object, len(removeOrder)+1)
	// find all objects and add them to objectsToDeleteByOrder
	for i := range installedPaks {
		for j := range installedPaks[i].Objects {
			md, err := installedPaks[i].Objects[j].Metadata()
			if err != nil {
				return err
			}
			order := stringInArrayIndex(removeOrder, md.Kind)
			if order == -1 {
				// set order to other objects
				order = len(objectsToDeleteByOrder) - 1
			}
			// add to corresponding order
			objectsToDeleteByOrder[order] = append(objectsToDeleteByOrder[order], installedPaks[i].Objects[j])
		}
	}

	for i := range objectsToDeleteByOrder {
		for j := range objectsToDeleteByOrder[i] {
			obj := objectsToDeleteByOrder[i][j]
			err := m.kubectl.Delete(namespace, obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
