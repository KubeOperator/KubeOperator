module github.com/fluxcd/helm-operator

go 1.13

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/fluxcd/flux v1.17.2-0.20200121140732-3903cf8e71c3
	github.com/fluxcd/helm-operator/pkg/install v0.0.0-00010101000000-000000000000
	github.com/ghodss/yaml v1.0.0
	github.com/go-kit/kit v0.9.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1
	github.com/gorilla/mux v1.7.1
	github.com/ncabatoff/go-seq v0.0.0-20180805175032-b08ef85ed833
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.24.0
	helm.sh/helm/v3 v3.2.3
	k8s.io/api v0.18.6 // kubernetes-1.16.2
	k8s.io/apiextensions-apiserver v0.18.6 // kubernetes-1.16.2
	k8s.io/apimachinery v0.18.6 // kubernetes-1.16.2
	k8s.io/cli-runtime v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/helm v2.16.1+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51 // kubernetes-1.16.2
)

// github.com/fluxcd/helm-operator/pkg/install lives in this very reprository, so use that
replace github.com/fluxcd/helm-operator/pkg/install => ./pkg/install

// Transitive requirement from Flux: remove when https://github.com/docker/distribution/pull/2905 is released.
replace github.com/docker/distribution => github.com/2opremio/distribution v0.0.0-20190419185413-6c9727e5e5de

// Transitive requirement from Helm.
replace github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0

// Force upgrade because of a transitive downgrade.
// github.com/fluxcd/helm-operator
// +-> github.com/fluxcd/flux@v1.15.0
//     +-> k8s.io/client-go@v11.0.0+incompatible
//     +-> github.com/fluxcd/helm-operator@v1.0.0-rc1
//         +-> k8s.io/client-go@v11.0.0+incompatible
//         +-> github.com/weaveworks/flux@v0.0.0-20190729133003-c78ccd3706b5
//             +-> k8s.io/client-go@v11.0.0+incompatible
replace k8s.io/client-go => k8s.io/client-go v0.18.6 // kubernetes-1.16.2

// Pin Flux to master branch to break weaveworks/flux circular dependency (to be removed on Flux 1.18)
replace github.com/fluxcd/flux => github.com/fluxcd/flux v1.17.2-0.20200121140732-3903cf8e71c3

// Patched release of Helm v3.0.3 until https://github.com/helm/helm/pull/7401 is merged
replace helm.sh/helm/v3 => github.com/hiddeco/helm/v3 v3.0.3-scheme-patched
