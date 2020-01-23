module github.com/mayadata-io/openebs-operator

go 1.13

require (
	github.com/ghodss/yaml v0.0.0-20150909031657-73d445a93680
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/pkg/errors v0.8.1
	k8s.io/api v0.0.0-20191005115622-2e41325d9e4b
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.0.0-20191008115822-1210218b4a26
	openebs.io/metac v0.1.1-0.20191209102635-9b94f129151c
)

replace openebs.io/metac => github.com/AmitKumarDas/metac v0.1.1-0.20191209102635-9b94f129151c
