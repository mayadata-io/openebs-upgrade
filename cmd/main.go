/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	"k8s.io/client-go/kubernetes"
	"openebs.io/metac/controller/generic"
	"openebs.io/metac/start"

	"github.com/golang/glog"
	"github.com/mayadata-io/openebs-operator/controller/openebs"
	"github.com/mayadata-io/openebs-operator/k8s"
)

// Command line flags
var (
	kubeconfig = flag.String(
		"kubeconfig", "",
		`Absolute path to the kubeconfig file.
		Required only when running outside the cluster.`,
	)
)

// main function is the entry point of this binary.
//
// This registers various controller (i.e. kubernetes reconciler)
// handler functions. Each handler function gets triggered due
// to any changes (add, update or delete) to configured watch
// resource.
func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Create the kubernetes client config.
	config, err := k8s.BuildConfig(*kubeconfig)
	if err != nil {
		glog.Error(err.Error())
		os.Exit(1)
	}
	// set the global variable Clientset value so that it
	// can be used globally.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Error(err.Error())
		os.Exit(1)
	}
	// set the global ClientSet variable so that it can be used globally.
	k8s.Clientset = clientset

	generic.AddToInlineRegistry("sync/openebs", openebs.Sync)

	start.Start()
}
