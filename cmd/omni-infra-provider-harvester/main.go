// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package main is the root cmd of the provider script.
package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	harvnetworkclient "github.com/harvester/harvester-network-controller/pkg/generated/clientset/versioned"
	harvclient "github.com/harvester/harvester/pkg/generated/clientset/versioned"
	"github.com/harvester/harvester/pkg/generated/clientset/versioned/scheme"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/siderolabs/omni/client/pkg/client"
	"github.com/siderolabs/omni/client/pkg/infra"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	kubeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	storageclient "k8s.io/client-go/kubernetes/typed/storage/v1"
	"k8s.io/client-go/rest"

	"github.com/ganawaj/omni-infra-provider-harvester/internal/pkg/provider"
	"github.com/ganawaj/omni-infra-provider-harvester/internal/pkg/provider/meta"
)

//go:embed data/schema.json
var schema string

//go:embed data/icon.svg
var icon []byte

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:          "provider",
	Short:        "Harvester Omni infrastructure provider",
	Long:         `Connects to Omni as an infra provider and manages VMs in Harvester`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		loggerConfig := zap.NewProductionConfig()

		logger, err := loggerConfig.Build(
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}

		baseConfig, err := kubeconfig.GetNonInteractiveClientConfig(cfg.kubeconfigFile).ClientConfig()
		if err != nil {
			return fmt.Errorf("failed to get client config: %w", err)
		}

		// fmt.Println("baseConfig", baseConfig)

		// Create a subresourced kubernetes rest client for harvester
		copyConfig := rest.CopyConfig(baseConfig)
		copyConfig.GroupVersion = &kubeschema.GroupVersion{Group: "subresources.kubevirt.io", Version: "v1"}
		copyConfig.APIPath = "/apis"
		copyConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

		restClient, err := rest.RESTClientFor(copyConfig)
		if err != nil {
			return fmt.Errorf("failed to get rest client: %w", err)
		}

		kubeClient, err := kubernetes.NewForConfig(baseConfig)
		if err != nil {
			return fmt.Errorf("failed to get kube client: %w", err)
		}
		storageClassClient, err := storageclient.NewForConfig(baseConfig)
		if err != nil {
			return fmt.Errorf("failed to get storage class client: %w", err)
		}
		harvClient, err := harvclient.NewForConfig(baseConfig)
		if err != nil {
			return fmt.Errorf("failed to get harvester client: %w", err)
		}
		harvNetworkClient, err := harvnetworkclient.NewForConfig(baseConfig)
		if err != nil {
			return fmt.Errorf("failed to get harvester network client: %w", err)
		}

		if cfg.omniAPIEndpoint == "" {
			return fmt.Errorf("omni-api-endpoint flag is not set")
		}

		harvesterClient := &provider.HarvesterClient{
			RestConfig:                baseConfig,
			KubeClient:                kubeClient,
			StorageClassClient:        storageClassClient,
			HarvesterClient:           harvClient,
			HarvesterNetworkClient:    harvNetworkClient,
			KubeVirtSubresourceClient: restClient,
		}

		provisioner := provider.NewProvisioner(harvesterClient, "")

		ip, err := infra.NewProvider(meta.ProviderID, provisioner, infra.ProviderConfig{
			Name:        cfg.providerName,
			Description: cfg.providerDescription,
			Icon:        base64.RawStdEncoding.EncodeToString(icon),
			Schema:      schema,
		})
		if err != nil {
			return fmt.Errorf("failed to create infra provider: %w", err)
		}

		logger.Info("starting infra provider")

		clientOptions := []client.Option{
			client.WithInsecureSkipTLSVerify(cfg.insecureSkipVerify),
		}

		if cfg.serviceAccountKey != "" {
			clientOptions = append(clientOptions, client.WithServiceAccount(cfg.serviceAccountKey))
		}

		return ip.Run(cmd.Context(), logger, infra.WithOmniEndpoint(cfg.omniAPIEndpoint), infra.WithClientOptions(
			clientOptions...,
		))
	},
}

var cfg struct {
	omniAPIEndpoint     string
	serviceAccountKey   string
	providerName        string
	providerDescription string
	kubeconfigFile      string
	dataVolumeMode      string
	insecureSkipVerify  bool
}

func main() {
	if err := app(); err != nil {
		os.Exit(1)
	}
}

func app() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.Flags().StringVar(&cfg.omniAPIEndpoint, "omni-api-endpoint", os.Getenv("OMNI_ENDPOINT"),
		"the endpoint of the Omni API, if not set, defaults to OMNI_ENDPOINT env var.")
	rootCmd.Flags().StringVar(&meta.ProviderID, "id", meta.ProviderID, "the id of the infra provider, it is used to match the resources with the infra provider label.")
	rootCmd.Flags().StringVar(&cfg.serviceAccountKey, "omni-service-account-key", os.Getenv("OMNI_SERVICE_ACCOUNT_KEY"), "Omni service account key, if not set, defaults to OMNI_SERVICE_ACCOUNT_KEY.")
	rootCmd.Flags().StringVar(&cfg.providerName, "provider-name", "Harvester", "provider name as it appears in Omni")
	rootCmd.Flags().StringVar(&cfg.providerDescription, "provider-description", "Harvester infrastructure provider", "Provider description as it appears in Omni")
	rootCmd.Flags().StringVar(&cfg.kubeconfigFile, "kubeconfig-file", "~/.kube/config", "Kubeconfig file to use to connect to the cluster where KubeVirt is running")
	rootCmd.Flags().BoolVar(&cfg.insecureSkipVerify, "insecure-skip-verify", false, "ignores untrusted certs on Omni side")
}
