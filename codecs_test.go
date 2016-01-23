package kupak

import (
	"testing"

	"github.com/ghodss/yaml"
)

var rc []byte = []byte(`apiVersion: v1
kind: ReplicationController
metadata:
  name: redis-standalone
spec:
  replicas: 1
  selector:
    name: redis-standalone
    mode: standalone
    provider: redis
  template:
    metadata:
      labels:
        name: redis-standalone
        mode: standalone
        provider: redis
        app: redis-standalone
    spec:
      containers:
      - name: redis-standalone
        image: kubernetes/redis:v1
        env:
        - name: MASTER
          value: "true"
        ports:
        - containerPort: 6060
        volumeMounts:
        - mountPath: /redis-master-data
          name: data
      volumes:
        - name: data
          emptyDir: {}`)

func TestObjectMetadata(t *testing.T) {
	obj := Object{
		data:         rc,
		marshaller:   yaml.Marshal,
		unmarshaller: yaml.Unmarshal,
	}
	meta, err := obj.Metadata()
	if err != nil {
		t.Log("Metadata Error", err)
		t.Fail()
	}
	if meta.Kind != "ReplicationController" {
		t.Fail()
	}
}

func TestReplicationController(t *testing.T) {
	obj, err := NewObject(rc, yaml.Marshal, yaml.Unmarshal)
	if err != nil {
		t.Log("RC Init Error", err)
		t.Fail()
	}
	rc := obj.ReplicationController()
	if rc.Spec.Template.Labels["mode"] != "standalone" {
		t.Log("RC Error")
		t.Fail()
	}
	obj.ReplicationController().Spec.Template.Labels["mode"] = "xxx"
	t.Log(obj.ReplicationController().Spec.Template.Labels["mode"])
	t.Fail()
}
