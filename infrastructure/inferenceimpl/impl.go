package inferenceimpl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"github.com/opensourceways/xihe-inference-evaluate/client"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/inference"
)

const MetaNameInference = "inference"

type CrdData struct {
	Group          string
	Version        string
	Name           string
	NameSpace      string
	Image          string
	GitlabEndPoint string
	XiheUser       string
	XiheUserToken  string
	ProjectName    string
	LastCommit     string
	ObsAk          string
	ObsSk          string
	ObsEndPoint    string
	ObsUtilPath    string
	ObsBucket      string
	ObsLfsPath     string
	StorageSize    int
	RecycleSeconds int
}

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

	res, err := impl.GetObj(impl.cfg, infer)

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

func (impl inferenceImpl) GetObj(cfg *Config, infer *domain.Inference) (*unstructured.Unstructured, error) {
	txtStr, err := ioutil.ReadFile("./crd-resource.yaml")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("crd").Parse(string(txtStr))
	if err != nil {
		return nil, err
	}

	var data = &CrdData{
		Group:          client.CrdGroup,
		Version:        client.CrdVersion,
		Name:           "",
		NameSpace:      "default",
		Image:          cfg.Image,
		GitlabEndPoint: cfg.GitlabEndpoint,
		XiheUser:       infer.User,
		XiheUserToken:  infer.UserToken,
		ProjectName:    infer.ProjectName.ProjectName(),
		LastCommit:     infer.LastCommit,
		ObsAk:          cfg.OBS.AccessKey,
		ObsSk:          cfg.OBS.SecretKey,
		ObsEndPoint:    cfg.OBS.Endpoint,
		ObsUtilPath:    cfg.OBS.OBSUtilPath,
		ObsBucket:      cfg.OBS.Bucket,
		ObsLfsPath:     cfg.OBS.LFSPath,
		StorageSize:    10,
		RecycleSeconds: 60,
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
