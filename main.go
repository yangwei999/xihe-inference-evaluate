package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-inference-evaluate/config"
	"github.com/opensourceways/xihe-inference-evaluate/controller"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/cloudimpl"
	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/inferenceimpl"
	"github.com/opensourceways/xihe-inference-evaluate/infrastructure/watchimpl"
	"github.com/opensourceways/xihe-inference-evaluate/k8sclient"
	"github.com/opensourceways/xihe-inference-evaluate/server"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)
	return o, err
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)

	if err != nil {
		logrus.Fatalf("new options failed, err:%s", err.Error())
	}

	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg := new(config.Config)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	if err := os.Remove(o.service.ConfigFile); err != nil {
		logrus.Fatalf("config file delete failed, err:%s", err.Error())
	}

	cli, err := k8sclient.Init(&cfg.K8sClient)
	if err != nil {
		logrus.Fatalf("k8s client init, err:%s", err.Error())
	}

	// inference
	inference, err := inferenceimpl.NewInference(&cli, &cfg.Inference, cfg.K8sClient)
	if err != nil {
		logrus.Fatalf("new inference service failed, err:%s", err.Error())
	}

	// cloud
	cloud, err := cloudimpl.NewCloud(&cli, &cfg.Cloud, cfg.K8sClient)
	if err != nil {
		logrus.Fatalf("new cloud service failed, err:%s", err.Error())
	}

	// controller
	controller.Init(log)

	// watcher
	w := watchimpl.NewWatcher(
		&cli,
		map[string]func(map[string]string, domain.ContainerDetail){
			inferenceimpl.MetaName(): inference.NotifyResult,
			cloudimpl.MetaName():     cloud.NotifyResult,
		},
	)

	w.Run()
	defer w.Exit()

	// run
	server.StartWebServer(
		o.service.Port, o.service.GracePeriod, cfg,
		&server.Service{
			Inference: inference,
			Cloud:     cloud,
		},
	)
}
