package cmd

import (
	"bytes"
	"io"

	k3dv1alpha5 "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	defaults "github.com/pthomison/k3auto/default"
	"github.com/pthomison/k3auto/internal/k3d"
	"github.com/sirupsen/logrus"
)

func checkError(err error) {
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func parseConfigFile(configPath string) (*k3dv1alpha5.SimpleConfig, error) {
	var clusterConfig *k3dv1alpha5.SimpleConfig
	var err error

	if configPath != "" {
		clusterConfig, err = k3d.ParseConfigFile(configPath, nil)
		if err != nil {
			return nil, err
		}
	} else {
		clusterConfig, err = k3d.ParseConfigFile(defaults.K3dConfigLocation, &defaults.K3dConfig)
		if err != nil {
			return nil, err
		}
	}

	return clusterConfig, nil
}

func yamlReadAndSplit(reader io.Reader) ([][]byte, error) {
	fb, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var objs [][]byte

	for _, obj := range bytes.Split(fb, []byte("---")) {
		if len(obj) != 0 {
			objs = append(objs, obj)
		}
	}

	return objs, err
}

// func k3autoDeploy() {
// 	deploymentFiles, err := defaults.DefaultDeployments.ReadDir(defaults.DefaultDeploymentsFolder)
// 	checkError(err)
// 	for _, v := range deploymentFiles {
// 		f, err := defaults.DefaultDeployments.Open(fmt.Sprintf("%v/%v", defaults.DefaultDeploymentsFolder, v.Name()))
// 		checkError(err)
// 		defer f.Close()

// 		fileObjects, err := yamlReadAndSplit(f)
// 		checkError(err)

// 		for _, obj := range fileObjects {
// 			obj, objType, err := k8s.ParseManifest(obj)
// 			checkError(err)

// 			logrus.Info("Deploying: ", objType)

// 			err = k8sC.Create(ctx, obj.(ctrlclient.Object))
// 			checkError(err)
// 		}
// 	}
// 	logrus.Info("Default Deployments Injected")
// }
