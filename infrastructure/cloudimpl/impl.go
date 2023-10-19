package cloudimpl

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	rpcclient "github.com/opensourceways/xihe-grpc-protocol/grpc/client"
	rpccloud "github.com/opensourceways/xihe-grpc-protocol/grpc/cloud"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/cloud"
	"github.com/opensourceways/xihe-inference-evaluate/k8sclient"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

const (
	keyId         = "id"
	metaNameCloud = "cloud"
)

func MetaName() string {
	return metaNameCloud
}

func NewCloud(
	cli *k8sclient.Client,
	cfg *Config,
	k8sConfig k8sclient.Config,
) (
	cloud.Cloud,
	error,
) {
	txtStr, err := ioutil.ReadFile(cfg.CRD.TemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("cloud").Parse(string(txtStr))
	if err != nil {
		return nil, err
	}

	rpcCli, err := rpcclient.NewCloudClient(cfg.RPCEndpoint)
	if err != nil {
		return nil, fmt.Errorf("new cloud rpc client error: %s", err.Error())
	}

	return &cloudImpl{
		cfg:         cfg,
		cli:         cli,
		rpcCli:      rpcCli,
		k8sConfig:   k8sConfig,
		crdTemplate: tmpl,
	}, nil
}

type cloudImpl struct {
	cfg         *Config
	cli         *k8sclient.Client
	rpcCli      *rpcclient.CloudClient
	k8sConfig   k8sclient.Config
	crdTemplate *template.Template
}

func (impl *cloudImpl) geneMetaName(index *domain.CloudPod) string {
	return fmt.Sprintf("%s-%s", metaNameCloud, index.PodId)
}

func (impl *cloudImpl) geneLabels(cloud *domain.CloudPod) map[string]string {
	return map[string]string{
		"type": metaNameCloud,
		keyId:  cloud.PodId,
	}
}

func (impl *cloudImpl) cloudIndexString(e *domain.CloudPod) string {
	return fmt.Sprintf(
		"Id:%s, meta name:%s",
		e.PodId, impl.geneMetaName(e),
	)
}

func (impl *cloudImpl) Create(cloud *domain.CloudPod) error {
	s := impl.cloudIndexString(cloud)
	logrus.Debugf("create cloud for %s.", s)

	res := new(unstructured.Unstructured)

	if err := impl.getObj(cloud, res); err != nil {
		return err
	}

	err := impl.cli.CreateCRD(res)

	logrus.Debugf(
		"create cloud:%s in %s, err:%v.",
		s, impl.k8sConfig.Namespace, err,
	)

	return err
}

func (impl *cloudImpl) NotifyResult(labels map[string]string, status domain.ContainerDetail) {
	cloud := rpccloud.CloudPod{
		Id: labels[keyId],
	}

	info := rpccloud.PodInfo{
		Error:     status.ErrorMsg,
		AccessURL: status.AccessUrl,
	}

	if err := impl.rpcCli.SetPodInfo(&cloud, &info); err != nil {
		logrus.Errorf("call cloud rpc error:%s", err.Error())
	} else {
		logrus.Debugf(
			"call rpc to set cloud(%s) info:(%s/%s)",
			cloud.Id,
			info.Error, info.AccessURL,
		)
	}
}

func (impl *cloudImpl) getObj(
	cloud *domain.CloudPod, obj *unstructured.Unstructured,
) error {
	crd := &impl.cfg.CRD
	k8sConfig := &impl.k8sConfig

	data := &crdData{
		Group:          k8sConfig.Group,
		Version:        k8sConfig.Version,
		CodeServer:     k8sConfig.Kind,
		Name:           impl.geneMetaName(cloud),
		NameSpace:      k8sConfig.Namespace,
		Image:          crd.CRDImage,
		User:           cloud.User,
		CPU:            crd.CRDCpuString(),
		Memory:         crd.CRDMemoryString(),
		StorageSize:    20,
		RecycleSeconds: int(cloud.SurvivalTime.SurvivalTime()),
		Labels:         impl.geneLabels(cloud),
		ContainerPort:  crd.CRDContainerPortString(),
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
	User           string
	CPU            string
	Memory         string
	StorageSize    int
	RecycleSeconds int
	Labels         map[string]string
	ContainerPort  string
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
