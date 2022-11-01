package inferenceimpl

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/opensourceways/xihe-inference-evaluate/client"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/inference"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

const MetaNameInference = "inference"

func NewInference(cfg *Config) inference.Inference {
	return inferenceImpl{}
}

type inferenceImpl struct {
}

func (impl inferenceImpl) Create(infer *domain.Inference) error {
	cli := client.GetDyna()
	resource, err, res := impl.GetResource()
	if err != nil {
		return err
	}

	res.Object["metadata"] = map[string]interface{}{
		"name":   impl.geneMetaName(&infer.InferenceIndex),
		"labels": impl.GeneLabels(infer),
	}

	dr := cli.Resource(resource).Namespace("default")
	_, err = dr.Create(context.TODO(), res, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (impl inferenceImpl) ExtendExpiry(infer *domain.InferenceIndex, expiry int64) error {
	cli := client.GetDyna()
	resource, err, _ := impl.GetResource()
	if err != nil {
		return err
	}

	get, err := cli.Resource(resource).Namespace("default").Get(context.TODO(), impl.geneMetaName(infer), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if sp, ok := get.Object["spec"]; ok {
		if spc, ok := sp.(map[string]interface{}); ok {
			spc["add"] = true
			spc["recycleAfterSeconds"] = expiry
		}
	}
	_, err = cli.Resource(resource).Namespace("default").Update(context.TODO(), get, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (impl inferenceImpl) GetResource() (schema.GroupVersionResource, error, *unstructured.Unstructured) {
	k, err, res := impl.resource()
	if err != nil {
		return schema.GroupVersionResource{}, err, nil
	}

	mapping, err := client.GetrestMapper().RESTMapping(k.GroupKind(), k.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err, nil
	}

	return mapping.Resource, nil, res
}

func (impl inferenceImpl) resource() (kind *schema.GroupVersionKind, err error, _ *unstructured.Unstructured) {
	var yamldata []byte
	yamldata, err = ioutil.ReadFile("crd-resource.yaml")
	if err != nil {
		return nil, err, nil
	}
	obj := &unstructured.Unstructured{}
	_, kind, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(yamldata, nil, obj)
	if err != nil {
		return nil, err, nil
	}
	return kind, nil, obj
}

func (impl inferenceImpl) geneMetaName(index *domain.InferenceIndex) string {
	return fmt.Sprintf("%s-%s", MetaNameInference, index.Id)
}

func (impl inferenceImpl) GeneLabels(infer *domain.Inference) map[string]string {
	m := make(map[string]string)
	m["id"] = infer.Id
	m["user"] = infer.Project.Owner.Account()
	m["project_id"] = infer.Project.Id
	m["last_commit"] = infer.LastCommit
	m["type"] = MetaNameInference
	return m
}
