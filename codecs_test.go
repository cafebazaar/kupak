package kupak

import "testing"

func TestObjectMetadata(t *testing.T) {
	obj := Object{data: rc}
	meta, err := obj.Metadata()
	if err != nil {
		t.Log("Metadata Error", err)
		t.Fail()
	}
	if meta.Kind != "ReplicationController" {
		t.Fail()
	}
}

func TestLabels(t *testing.T) {
	obj := Object{data: rc}
	metadata, err := obj.Metadata()
	if err != nil {
		t.Log("can't get metadata - ", err)
		t.Fail()
	}
	labels := metadata.Labels

	labels["hi"] = "hello"
	if err := obj.SetLabels(labels); err != nil {
		t.Log("can't set labels - ", err)
		t.Fail()
	}

	metadata, err = obj.Metadata()
	if err != nil {
		t.Log("can't get metadata - ", err)
		t.Fail()
	}
	labels = metadata.Labels

	if labels["hi"] != "hello" {
		t.Log("labels doesn't changed")
		t.Fail()
	}
	if labels["test"] != "hi" {
		t.Log("labels are corrupted")
		t.Fail()
	}
}

func TestTemplateMetadata(t *testing.T) {
	obj := Object{data: rc}
	meta, err := obj.TemplateMetadata()
	if err != nil {
		t.Log("metadata error", err)
		t.Fail()
	}
	if meta.Labels["name"] != "redis-standalone" {
		t.Log("template metadata error")
		t.Fail()
	}
}

func TestTemplateLabels(t *testing.T) {
	obj := Object{data: rc}
	metadata, err := obj.TemplateMetadata()
	if err != nil {
		t.Log("can't get metadata - ", err)
		t.Fail()
	}
	labels := metadata.Labels
	labels["hi"] = "hello"
	if err := obj.SetTemplateLabels(labels); err != nil {
		t.Log("can't set labels - ", err)
		t.Fail()
	}
	if labels["hi"] != "hello" {
		t.Log("labels doesn't changed")
		t.Fail()
	}
	if labels["name"] != "redis-standalone" {
		t.Log("labels are corrupted")
		t.Fail()
	}
}

var rc []byte = []byte(`apiVersion: v1
kind: ReplicationController
metadata:
  name: redis-standalone
  labels:
    test: hi
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
