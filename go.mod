module github.com/lsst/qserv-operator

go 1.13

require (
	github.com/go-logr/logr v0.3.0
	github.com/go-openapi/spec v0.19.3
	github.com/go-test/deep v1.0.7
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/iancoleman/strcase v0.1.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	sigs.k8s.io/controller-runtime v0.7.2
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)
