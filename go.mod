module github.com/23technologies/gardener-extension-runtime-kata

go 1.16

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/ahmetb/gen-crd-api-reference-docs v0.2.0
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/gardener/gardener v1.43.1
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220325170049-de3da57026de // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/tools v0.1.9
	k8s.io/apimachinery v0.23.5
	k8s.io/code-generator v0.23.5
	k8s.io/component-base v0.23.5
	k8s.io/klog/v2 v2.50.0 // indirect
	k8s.io/metrics v0.23.5 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/controller-runtime v0.11.2
)

replace (
	// github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	k8s.io/api => k8s.io/api v0.22.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.2
	k8s.io/apiserver => k8s.io/apiserver v0.22.2
	k8s.io/client-go => k8s.io/client-go v0.22.2
	k8s.io/code-generator => k8s.io/code-generator v0.22.2
	k8s.io/component-base => k8s.io/component-base v0.22.2
	k8s.io/helm => k8s.io/helm v2.13.1+incompatible
// sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.10.2
)
