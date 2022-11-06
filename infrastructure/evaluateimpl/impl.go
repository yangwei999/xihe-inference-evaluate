package evaluateimpl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/evaluate"
	"github.com/opensourceways/xihe-inference-evaluate/k8sclient"
)

const metaNameEvaluate = "evaluate"

func MetaName() string {
	return metaNameEvaluate
}

func NewEvaluate(cfg *Config) (evaluate.Evaluate, error) {
	txtStr, err := ioutil.ReadFile(cfg.CRD.TemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("evaluate").Parse(string(txtStr))
	if err != nil {
		return nil, err
	}

	return &evaluateImpl{
		cfg:         cfg,
		crdTemplate: tmpl,
	}, nil
}

type evaluateImpl struct {
	cfg         *Config
	crdTemplate *template.Template
}

func (impl *evaluateImpl) evaluateIndexString(e *domain.EvaluateIndex) string {
	return fmt.Sprintf(
		"%s/%s/%s/%s, meta name:%s",
		e.Project.Owner.Account(), e.Project.Id,
		e.TrainingId, e.Id, impl.geneMetaName(e),
	)
}

func (impl *evaluateImpl) CreateCustom(ce *domain.CustomEvaluate) error {
	s := impl.evaluateIndexString(&ce.EvaluateIndex)
	logrus.Debugf("receive custom evaluate for %s.", s)

	cli := k8sclient.GetDyna()
	resource := k8sclient.GetResource()

	res := new(unstructured.Unstructured)

	if err := impl.getCustomObj(ce, res); err != nil {
		return err
	}

	dr := cli.Resource(resource).Namespace(impl.cfg.CRD.CRDNamespace)

	_, err := dr.Create(context.TODO(), res, metav1.CreateOptions{})

	logrus.Debugf(
		"gen crd for custom evaluate:%s in %s/%s, err:%v.",
		s, resource, impl.cfg.CRD.CRDNamespace, err,
	)

	return err
}

func (impl *evaluateImpl) CreateStandard(se *domain.StandardEvaluate) error {
	s := impl.evaluateIndexString(&se.EvaluateIndex)
	logrus.Debugf("receive standard evaluate for %s.", s)

	cli := k8sclient.GetDyna()
	resource := k8sclient.GetResource()

	res := new(unstructured.Unstructured)

	if err := impl.getStandardObj(se, res); err != nil {
		return err
	}

	dr := cli.Resource(resource).Namespace(impl.cfg.CRD.CRDNamespace)

	_, err := dr.Create(context.TODO(), res, metav1.CreateOptions{})

	logrus.Debugf(
		"gen crd for standard evaluate:%s in %s/%s, err:%v.",
		s, resource, impl.cfg.CRD.CRDNamespace, err,
	)

	return err
}

func (impl *evaluateImpl) geneMetaName(eva *domain.EvaluateIndex) string {
	return fmt.Sprintf("%s-%s", metaNameEvaluate, eva.Id)
}

func (impl *evaluateImpl) geneLabels(eva *domain.EvaluateIndex) map[string]string {
	m := make(map[string]string)
	m["id"] = eva.Id
	m["user"] = eva.Project.Owner.Account()
	m["project_id"] = eva.Project.Id
	m["training_id"] = eva.TrainingId
	m["type"] = metaNameEvaluate
	return m
}

func (impl *evaluateImpl) getCustomObj(
	ceva *domain.CustomEvaluate,
	obj *unstructured.Unstructured,
) error {
	name := impl.geneMetaName(&ceva.EvaluateIndex)
	labels := impl.geneLabels(&ceva.EvaluateIndex)
	crd := &impl.cfg.CRD

	data := crdData{
		Group:          k8sclient.Cfg.Group,
		Version:        k8sclient.Cfg.Version,
		CodeServer:     k8sclient.Cfg.Kind,
		Name:           name,
		NameSpace:      crd.CRDNamespace,
		Image:          crd.CRDImage,
		ObsAk:          impl.cfg.OBS.AccessKey,
		ObsSk:          impl.cfg.OBS.SecretKey,
		ObsEndPoint:    impl.cfg.OBS.Endpoint,
		ObsUtilPath:    impl.cfg.OBS.OBSUtilPath,
		ObsBucket:      impl.cfg.OBS.Bucket,
		StorageSize:    10,
		RecycleSeconds: ceva.SurvivalTime,
		CPU:            crd.CRDCpuString(),
		Memory:         crd.CRDMemoryString(),
		Labels:         labels,
		OBSPath:        ceva.AimPath,
		EvaluateType:   ceva.Type(),
	}

	return data.genTemplate(impl.crdTemplate, obj)
}

func (impl *evaluateImpl) getStandardObj(
	seva *domain.StandardEvaluate,
	obj *unstructured.Unstructured,
) error {
	name := impl.geneMetaName(&seva.EvaluateIndex)
	labels := impl.geneLabels(&seva.EvaluateIndex)
	crd := &impl.cfg.CRD

	data := crdData{
		Group:          k8sclient.Cfg.Group,
		Version:        k8sclient.Cfg.Version,
		CodeServer:     k8sclient.Cfg.Kind,
		Name:           name,
		NameSpace:      crd.CRDNamespace,
		Image:          crd.CRDImage,
		ObsAk:          impl.cfg.OBS.AccessKey,
		ObsSk:          impl.cfg.OBS.SecretKey,
		ObsEndPoint:    impl.cfg.OBS.Endpoint,
		ObsUtilPath:    impl.cfg.OBS.OBSUtilPath,
		ObsBucket:      impl.cfg.OBS.Bucket,
		StorageSize:    10,
		RecycleSeconds: seva.SurvivalTime,
		CPU:            crd.CRDCpuString(),
		Memory:         crd.CRDMemoryString(),
		Labels:         labels,
		OBSPath:        seva.LogPath,
		EvaluateType:   seva.Type(),
		LearningScope:  seva.LearningRateScope.String(),
		BatchScope:     seva.BatchSizeScope.String(),
		MomentumScope:  seva.MomentumScope.String(),
	}

	return data.genTemplate(impl.crdTemplate, obj)
}

type crdData struct {
	Group          string
	Version        string
	CodeServer     string
	Name           string
	NameSpace      string
	Image          string
	CPU            string
	Memory         string
	StorageSize    int
	RecycleSeconds int
	Labels         map[string]string

	ObsAk         string
	ObsSk         string
	ObsEndPoint   string
	ObsUtilPath   string
	ObsBucket     string
	OBSPath       string
	EvaluateType  string
	LearningScope string
	BatchScope    string
	MomentumScope string
}

func (data *crdData) genTemplate(tmpl *template.Template, obj *unstructured.Unstructured) error {
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}

	_, _, err := yaml.NewDecodingSerializer(
		unstructured.UnstructuredJSONScheme,
	).Decode(buf.Bytes(), nil, obj)

	return err
}
