package evaluateimpl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opensourceways/xihe-inference-evaluate/client"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/evaluate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

const MetaNameEvaluate = "evaluate"

func NewEvaluate(cfg *Config) evaluate.Evaluate {
	return &evaluateImpl{
		cfg: cfg,
	}
}

type evaluateImpl struct {
	cfg *Config
}

func (impl *evaluateImpl) CreateCustom(eva *domain.CustomEvaluate) error {
	cli := client.GetDyna()
	resource := client.GetResource2()

	res, err := impl.GetObj(eva)
	dr := cli.Resource(resource).Namespace("default")

	_, err = dr.Create(context.TODO(), res, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (impl *evaluateImpl) CreateStandard(*domain.StandardEvaluate) error {

	return nil

}

func (impl *evaluateImpl) geneMetaName(eva *domain.CustomEvaluate) string {
	return fmt.Sprintf("%s-%s", MetaNameEvaluate, eva.Id)
}

func (impl *evaluateImpl) GeneLabels(eva *domain.CustomEvaluate) map[string]string {
	m := make(map[string]string)
	m["id"] = eva.Id
	m["user"] = eva.Project.Owner.Account()
	m["project_id"] = eva.Project.Id
	m["training_id"] = eva.TrainingId
	m["type"] = MetaNameEvaluate
	return m
}

func (impl *evaluateImpl) GetObj(eva *domain.CustomEvaluate) (*unstructured.Unstructured, error) {
	name := impl.geneMetaName(eva)
	labels := impl.GeneLabels(eva)

	tmpl, err := client.GetTemplate()
	if err != nil {
		return nil, err
	}

	var data = &client.CrdData{
		Group:          client.CrdGroup,
		Version:        client.CrdVersion,
		Name:           name,
		NameSpace:      "default",
		Image:          impl.cfg.Image,
		ObsAk:          impl.cfg.OBS.AccessKey,
		ObsSk:          impl.cfg.OBS.SecretKey,
		ObsEndPoint:    impl.cfg.OBS.Endpoint,
		ObsUtilPath:    impl.cfg.OBS.OBSUtilPath,
		ObsBucket:      impl.cfg.OBS.Bucket,
		ObsLfsPath:     impl.cfg.OBS.LFSPath,
		StorageSize:    10,
		RecycleSeconds: eva.SurvivalTime,
		Labels:         labels,
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	_, _, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(buf.Bytes(), nil, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
