package inferenceimpl

import (
	"bytes"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"github.com/opensourceways/xihe-inference-evaluate/client"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/inference"
)

const MetaNameInference = "inference"

func NewInference(cfg *Config) inference.Inference {
	return inferenceImpl{
		cfg: cfg,
	}
}

type inferenceImpl struct {
	cfg *Config
}

func (impl inferenceImpl) Create(infer *domain.Inference) error {
	cli := client.GetDyna()
	resource := client.GetResource2()

	res, err := impl.GetObj(infer)

	dr := cli.Resource(resource).Namespace("default")
	_, err = dr.Create(context.TODO(), res, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (impl inferenceImpl) ExtendSurvivalTime(infer *domain.InferenceIndex, timeToExtend int) error {
	cli := client.GetDyna()
	resource := client.GetResource2()

	get, err := cli.Resource(resource).Namespace("default").Get(context.TODO(), impl.geneMetaName(infer), metav1.GetOptions{})
	if err != nil {
		return err
	}

	if sp, ok := get.Object["spec"]; ok {
		if spc, ok := sp.(map[string]interface{}); ok {
			spc["increaseRecycleSeconds"] = true
			spc["recycleAfterSeconds"] = timeToExtend
		}
	}
	_, err = cli.Resource(resource).Namespace("default").Update(context.TODO(), get, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
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

func (impl inferenceImpl) GetObj(infer *domain.Inference) (*unstructured.Unstructured, error) {
	name := impl.geneMetaName(&infer.InferenceIndex)
	labels := impl.GeneLabels(infer)

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
		GitlabEndPoint: impl.cfg.GitlabEndpoint,
		XiheUser:       infer.User,
		XiheUserToken:  infer.UserToken,
		ProjectName:    infer.ProjectName.ProjectName(),
		LastCommit:     infer.LastCommit,
		ObsAk:          impl.cfg.OBS.AccessKey,
		ObsSk:          impl.cfg.OBS.SecretKey,
		ObsEndPoint:    impl.cfg.OBS.Endpoint,
		ObsUtilPath:    impl.cfg.OBS.OBSUtilPath,
		ObsBucket:      impl.cfg.OBS.Bucket,
		ObsLfsPath:     impl.cfg.OBS.LFSPath,
		StorageSize:    10,
		RecycleSeconds: infer.SurvivalTime,
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
