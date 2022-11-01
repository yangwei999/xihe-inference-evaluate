package client

import (
	"io/ioutil"
	"log"
	"os/user"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sConfig *rest.Config
	k8sClient *kubernetes.Clientset
	dyna      dynamic.Interface
	restm     *restmapper.DeferredDiscoveryRESTMapper
)

func getHome() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}

	return u.HomeDir
}

func Init() (err error) {
	k8sConfig, err = clientcmd.BuildConfigFromFlags("", getHome()+"/.kube/config")
	if err != nil {
		log.Println(err)
		return
	}

	k8sClient, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Println(err)
		return
	}
	dyna, err = dynamic.NewForConfig(k8sConfig)
	if err != nil {
		log.Println(err)
		return
	}

	dis, err := discovery.NewDiscoveryClientForConfig(k8sConfig)
	if err != nil {
		log.Println("NewDiscoveryClientForConfig err", err)
		return
	}

	restm = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dis))
	return nil
}

func GetClient() *kubernetes.Clientset {
	return k8sClient
}

func GetDyna() dynamic.Interface {
	return dyna
}

func GetrestMapper() *restmapper.DeferredDiscoveryRESTMapper {
	return restm
}

func GetK8sConfig() *rest.Config {
	return k8sConfig
}

func GetResource() (schema.GroupVersionResource, error, *unstructured.Unstructured) {
	k, err, res := Resource()
	if err != nil {
		return schema.GroupVersionResource{}, err, nil
	}

	mapping, err := GetrestMapper().RESTMapping(k.GroupKind(), k.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err, nil
	}

	return mapping.Resource, nil, res
}

func Resource() (kind *schema.GroupVersionKind, err error, _ *unstructured.Unstructured) {
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
