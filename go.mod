module github.com/ganawaj/omni-infra-provider-harvester

go 1.24.1

toolchain go1.24.3

replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	github.com/rancher/lasso => github.com/rancher/lasso v0.0.0-20241202185148-04649f379358
	github.com/rancher/rancher => github.com/rancher/rancher v0.0.0-20240919204204-3da2ae0cabd1
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20240919204204-3da2ae0cabd1
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20240919204204-3da2ae0cabd1
	go.qase.io/client => github.com/rancher/qase-go/client v0.0.0-20231114201952-65195ec001fa
	gopkg.in/yaml.v3 => github.com/unix4ever/yaml v0.0.0-20220527175918-f17b0f05cf2c
	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.15.1-rancher2
	k8s.io/api => k8s.io/api v0.31.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.31.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.31.5
	k8s.io/apiserver => k8s.io/apiserver v0.31.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.31.5
	k8s.io/client-go => k8s.io/client-go v0.31.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.31.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.31.5
	k8s.io/code-generator => k8s.io/code-generator v0.31.5
	k8s.io/component-base => k8s.io/component-base v0.31.5
	k8s.io/component-helpers => k8s.io/component-helpers v0.31.5
	k8s.io/controller-manager => k8s.io/controller-manager v0.31.5
	k8s.io/cri-api => k8s.io/cri-api v0.31.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.31.5
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.31.5
	k8s.io/endpointslice => k8s.io/endpointslice v0.31.5
	k8s.io/kms => k8s.io/kms v0.31.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.31.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.31.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.31.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.31.5
	k8s.io/kubectl => k8s.io/kubectl v0.31.5
	k8s.io/kubelet => k8s.io/kubelet v0.31.5
	k8s.io/kubernetes => k8s.io/kubernetes v1.31.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.31.5
	k8s.io/metrics => k8s.io/metrics v0.31.5
	k8s.io/mount-utils => k8s.io/mount-utils v0.31.5
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.31.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.31.5
)

require (
	github.com/cosi-project/runtime v0.10.5
	github.com/google/uuid v1.6.0
	github.com/harvester/harvester v1.5.0
	github.com/harvester/harvester-network-controller v0.3.1
	github.com/planetscale/vtprotobuf v0.6.1-0.20241121165744-79df5c4772f2
	github.com/rancher/wrangler/v3 v3.1.0
	github.com/siderolabs/go-pointer v1.0.1
	github.com/siderolabs/omni/client v0.50.0
	github.com/spf13/cobra v1.9.1
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.36.6
	k8s.io/api v0.33.0
	k8s.io/apimachinery v0.34.0-alpha.0
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/api v1.5.0
)

require (
	cel.dev/expr v0.23.1 // indirect
	emperror.dev/errors v0.8.1 // indirect
	github.com/ProtonMail/go-crypto v1.2.0 // indirect
	github.com/ProtonMail/go-mime v0.0.0-20230322103455-7d82a3887f2f // indirect
	github.com/ProtonMail/gopenpgp/v2 v2.8.3 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cisco-open/operator-tools v0.34.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/containerd/go-cni v1.1.12 // indirect
	github.com/containernetworking/cni v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/iancoleman/orderedmap v0.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/jsimonetti/rtnetlink/v2 v2.0.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.7.0 // indirect
	github.com/k8snetworkplumbingwg/whereabouts v0.8.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kube-logging/logging-operator/pkg/sdk v0.11.1-0.20240314152935-421fefebc813 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/longhorn/longhorn-manager v1.8.1 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mdlayher/ethtool v0.4.0 // indirect
	github.com/mdlayher/genetlink v1.3.2 // indirect
	github.com/mdlayher/netlink v1.7.3-0.20250113171957-fbb4dce95f42 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/runtime-spec v1.2.1 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/petermattis/goid v0.0.0-20250319124200-ccd6737f222a // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.72.0 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.0 // indirect
	github.com/rancher/aks-operator v1.9.2 // indirect
	github.com/rancher/eks-operator v1.9.2 // indirect
	github.com/rancher/fleet/pkg/apis v0.10.0 // indirect
	github.com/rancher/gke-operator v1.9.2 // indirect
	github.com/rancher/lasso v0.0.0-20241202185148-04649f379358 // indirect
	github.com/rancher/norman v0.0.0-20240708202514-a0127673d1b9 // indirect
	github.com/rancher/rancher/pkg/apis v0.0.0 // indirect
	github.com/rancher/rke v1.6.2 // indirect
	github.com/rancher/system-upgrade-controller/pkg/apis v0.0.0-20240301001845-4eacc2dabbde // indirect
	github.com/rancher/wrangler v1.1.2 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/siderolabs/crypto v0.5.1 // indirect
	github.com/siderolabs/gen v0.8.1 // indirect
	github.com/siderolabs/go-api-signature v0.3.6 // indirect
	github.com/siderolabs/image-factory v0.7.0 // indirect
	github.com/siderolabs/net v0.4.0 // indirect
	github.com/siderolabs/proto-codec v0.1.2 // indirect
	github.com/siderolabs/protoenc v0.2.2 // indirect
	github.com/siderolabs/talos/pkg/machinery v1.10.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/grpc v1.72.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.3 // indirect
	k8s.io/apiextensions-apiserver v0.32.3 // indirect
	k8s.io/apiserver v0.32.3 // indirect
	k8s.io/component-base v0.32.3 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.31.5 // indirect
	k8s.io/kubernetes v1.32.2 // indirect
	k8s.io/utils v0.0.0-20250502105355-0f33e8f1c979 // indirect
	kubevirt.io/containerized-data-importer-api v1.61.2 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.2.4 // indirect
	sigs.k8s.io/cli-utils v0.37.2 // indirect
	sigs.k8s.io/cluster-api v1.7.3 // indirect
	sigs.k8s.io/controller-runtime v0.20.3 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.7.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
