package kupak

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestObjectMetadata(t *testing.T) {
	obj := Object{
		data: `
apiVersion: v1
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
        - containerPort: $(.port)
        volumeMounts:
        - mountPath: /redis-master-data
          name: data
      volumes:
        - name: data
          emptyDir: {}
		`,
		marshaller:   yaml.Marshal,
		unmarshaller: yaml.Unmarshal,
	}
	meta, err := obj.Metadata()
	if err != nil {
		t.Fail()
	}
	if meta.Kind != "ReplicationController" {
		t.Fail()
	}
}
