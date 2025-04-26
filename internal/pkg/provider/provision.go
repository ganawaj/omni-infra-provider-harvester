// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package provider implements KubeVirt infra provider core.
package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	harvnetworkclient "github.com/harvester/harvester-network-controller/pkg/generated/clientset/versioned"
	"github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	harvclient "github.com/harvester/harvester/pkg/generated/clientset/versioned"
	"github.com/siderolabs/go-pointer"
	"github.com/siderolabs/omni/client/pkg/constants"
	"github.com/siderolabs/omni/client/pkg/infra/provision"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	storageclient "k8s.io/client-go/kubernetes/typed/storage/v1"
	"k8s.io/client-go/rest"
	kvv1 "kubevirt.io/api/core/v1"

	"github.com/ganawaj/omni-infra-provider-harvester/internal/pkg/provider/resources"
)

// HarvesterClient is a client for the Harvester API.
type HarvesterClient struct {
	RestConfig                *rest.Config
	KubeVirtSubresourceClient *rest.RESTClient
	KubeClient                *kubernetes.Clientset
	StorageClassClient        *storageclient.StorageV1Client
	HarvesterClient           *harvclient.Clientset
	HarvesterNetworkClient    *harvnetworkclient.Clientset
}

// Provisioner implements Talos emulator infra provider.
type Provisioner struct {
	harvesterClient *HarvesterClient
	namespace       string
}

// NewProvisioner creates a new provisioner.
func NewProvisioner(harvesterClient *HarvesterClient, namespace string) *Provisioner {
	return &Provisioner{
		harvesterClient: harvesterClient,
		namespace:       namespace,
	}
}

// ProvisionSteps implements infra.Provisioner.
//
//nolint:gocognit,gocyclo,cyclop,maintidx
func (p *Provisioner) ProvisionSteps() []provision.Step[*resources.Machine] {
	return []provision.Step[*resources.Machine]{
		// Validate the request
		provision.NewStep("validateRequest", func(_ context.Context, _ *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			if len(pctx.GetRequestID()) > 62 {
				return fmt.Errorf("the machine request name can not be longer than 63 characters")
			}

			return nil
		}),

		provision.NewStep("namespace", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			var data Data

			err := pctx.UnmarshalProviderData(&data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal provider data: %w", err)
			}

			if data.Namespace == "" {
				logger.Error("namespace is not set")

				return fmt.Errorf("namespace is not set")
			}

			nsClient := p.harvesterClient.KubeClient.CoreV1().Namespaces()

			ns, err := nsClient.Get(ctx, data.Namespace, k8smetav1.GetOptions{})
			if errors.IsNotFound(err) {
				logger.Error("namespace not found", zap.String("namespace", data.Namespace))

				return provision.NewRetryInterval(time.Second * 10)
			}

			if err != nil {
				logger.Error("failed to get the namespace", zap.Error(err))

				return provision.NewRetryInterval(time.Second * 10)
			}

			p.namespace = ns.Name

			return nil
		}),

		// Create the schematic
		provision.NewStep("createSchematic", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			schematic, err := pctx.GenerateSchematicID(ctx, logger,
				provision.WithExtraKernelArgs("console=ttyS0,38400n8"),
			)
			if err != nil {
				return err
			}

			pctx.State.TypedSpec().Value.Schematic = schematic

			return nil
		}),

		// Ensure the volume
		provision.NewStep("ensureVolume", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			pctx.State.TypedSpec().Value.TalosVersion = pctx.GetTalosVersion()

			url, err := url.Parse(constants.ImageFactoryBaseURL)
			if err != nil {
				logger.Error("failed to parse image factory base URL", zap.Error(err))

				return err
			}

			var data Data

			err = pctx.UnmarshalProviderData(&data)
			if err != nil {
				logger.Error("failed to unmarshal provider data", zap.Error(err))

				return err
			}

			url = url.JoinPath("image",
				pctx.State.TypedSpec().Value.Schematic,
				pctx.GetTalosVersion(),
				fmt.Sprintf("nocloud-%s.qcow2", data.Architecture),
			)

			hash := sha256.New()

			if _, err = hash.Write([]byte(url.String())); err != nil {
				return err
			}

			volumeID := hex.EncodeToString(hash.Sum(nil))
			volumeName := fmt.Sprintf("talos-%s", volumeID)
			volumeIdentifier := volumeName[:16]

			found, err := p.harvesterClient.HarvesterClient.HarvesterhciV1beta1().
				VirtualMachineImages(p.namespace).List(ctx, k8smetav1.ListOptions{
				LabelSelector: "omni.siderolabs.io/volume-id=" + volumeIdentifier,
			})
			if err != nil {
				logger.Error("failed to list the base talos image", zap.Error(err))

				return err
			}

			if len(found.Items) > 0 {
				logger.Info("base talos image already exists, skipping creation", zap.String("volumeName", volumeName))
				pctx.State.TypedSpec().Value.VolumeId = found.Items[0].ObjectMeta.Name

				return nil
			}

			logger.Info("base talos image not found, creating it", zap.String("volumeName", volumeName))

			// Validate the storage class
			_, err = p.harvesterClient.StorageClassClient.
				StorageClasses().Get(ctx, data.StorageClass, k8smetav1.GetOptions{})
			if err != nil {
				logger.Error("failed to get the storage class", zap.Error(err))

				return err
			}

			// Create the base talos image
			vmImage := &v1beta1.VirtualMachineImage{
				ObjectMeta: k8smetav1.ObjectMeta{
					GenerateName: fmt.Sprintf("%s-", volumeName),
					Namespace:    p.namespace,
					Labels: map[string]string{
						"tag.harvesterhci.io/created-by": "omni-infra-provider-harvester",
						"tag.harvesterhci.io/managed-by": "omni",
						"harvesterhci.io/creator":        "omni-infra-provider-harvester",
						"omni.siderolabs.io/volume-id":   volumeIdentifier,
					},
					Annotations: map[string]string{
						"harvesterhci.io/storageClassName": data.StorageClass,
					},
				},
				Spec: v1beta1.VirtualMachineImageSpec{
					DisplayName: pctx.GetRequestID(),
					SourceType:  v1beta1.VirtualMachineImageSourceTypeDownload,
					URL:         url.String(),
					Retry:       3,
				},
			}

			// Create the Image
			image, err := p.harvesterClient.HarvesterClient.HarvesterhciV1beta1().
				VirtualMachineImages(p.namespace).Create(ctx, vmImage, k8smetav1.CreateOptions{})
			if err != nil {
				logger.Error("failed to create the base talos image", zap.Error(err))

				return provision.NewRetryInterval(time.Second * 10)
			}

			watch, err := p.harvesterClient.HarvesterClient.
				HarvesterhciV1beta1().
				VirtualMachineImages(p.namespace).
				Watch(ctx, k8smetav1.ListOptions{
					FieldSelector:  fmt.Sprintf("metadata.name=%s", image.Name),
					TimeoutSeconds: pointer.To(int64(60)),
				})
			if err != nil {
				logger.Error("failed to watch the base talos image", zap.Error(err))

				return provision.NewRetryInterval(time.Second * 10)
			}

			defer watch.Stop()

			var final *v1beta1.VirtualMachineImage
			for event := range watch.ResultChan() {
				if watchImage, ok := event.Object.(*v1beta1.VirtualMachineImage); ok {
					if watchImage.Status.Progress == 100 {
						final = watchImage

						break
					}
					logger.Info("base talos image creation in progress", zap.String("volumeName", volumeName), zap.Int("progress", watchImage.Status.Progress))
				}
			}

			watch.Stop()

			if final != nil {
				logger.Info("base talos image creation completed", zap.String("volumeName", volumeName), zap.Int("progress", final.Status.Progress))
				pctx.State.TypedSpec().Value.VolumeId = final.ObjectMeta.Name
				logger.Info("base talos image created", zap.String("volumeName", volumeName))

				return nil
			}

			logger.Error("base talos image creation failed", zap.String("volumeName", volumeName))

			return provision.NewRetryInterval(time.Second * 10)
		}),

		// Create the machine
		provision.NewStep("createPVC", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			if pctx.State.TypedSpec().Value.Uuid == "" {
				pctx.State.TypedSpec().Value.Uuid = uuid.NewString()
			}

			var data Data

			err := pctx.UnmarshalProviderData(&data)
			if err != nil {
				logger.Error("failed to unmarshal provider data", zap.Error(err))

				return err
			}

			pvcName := fmt.Sprintf("%s-%s-%s", pctx.GetRequestID(), "disk-0", pctx.State.TypedSpec().Value.Uuid[0:8])

			_, err = p.harvesterClient.KubeClient.CoreV1().
				PersistentVolumeClaims(p.namespace).Get(ctx, pvcName, k8smetav1.GetOptions{})
			if err == nil {
				logger.Info("PVC already exists", zap.String("pvcName", pvcName))

				return nil
			}
			if !errors.IsNotFound(err) {
				logger.Error("failed to get the PVC", zap.Error(err))

				return err
			}

			pvc := &v1.PersistentVolumeClaim{
				ObjectMeta: k8smetav1.ObjectMeta{
					Name:      pvcName,
					Namespace: p.namespace,
					Labels: map[string]string{
						"omni.siderolabs.io/volume-id": pctx.State.TypedSpec().Value.VolumeId,
					},
					Annotations: map[string]string{
						"harvesterhci.io/imageId": fmt.Sprintf("%s/%s", p.namespace, pctx.State.TypedSpec().Value.VolumeId),
					},
				},
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						v1.ReadWriteMany,
					},
					Resources: v1.VolumeResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse(fmt.Sprintf("%dGi", data.DiskSize)),
						},
					},
					VolumeMode:       pointer.To(v1.PersistentVolumeBlock),
					StorageClassName: pointer.To("longhorn-" + pctx.State.TypedSpec().Value.VolumeId),
				},
			}

			_, err = p.harvesterClient.KubeClient.CoreV1().PersistentVolumeClaims(p.namespace).Create(ctx, pvc, k8smetav1.CreateOptions{})
			if err != nil {
				logger.Error("failed to create the PVC", zap.Error(err))

				return err
			}

			watch, err := p.harvesterClient.KubeClient.CoreV1().PersistentVolumeClaims(p.namespace).
				Watch(ctx, k8smetav1.ListOptions{
					FieldSelector:  fmt.Sprintf("metadata.name=%s", pvcName),
					TimeoutSeconds: pointer.To(int64(60)),
				})
			if err != nil {
				logger.Error("failed to watch the PVC", zap.Error(err))

				return err
			}

			defer watch.Stop()

			for event := range watch.ResultChan() {
				if watchPVC, ok := event.Object.(*v1.PersistentVolumeClaim); ok {
					if watchPVC.Status.Phase == v1.ClaimBound {
						logger.Info("PVC creation completed", zap.String("pvcName", pvcName), zap.String("phase", string(watchPVC.Status.Phase)))

						return nil
					}
					logger.Info("PVC creation in progress", zap.String("pvcName", pvcName), zap.String("phase", string(watchPVC.Status.Phase)))
				}
			}

			return provision.NewRetryInterval(time.Second * 10)
		}),

		// Create the machine
		provision.NewStep("syncMachine", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			logger = logger.With(zap.String("id", pctx.State.TypedSpec().Value.Uuid))

			defer func() {
				pctx.SetMachineUUID(pctx.State.TypedSpec().Value.Uuid)
			}()

			var data Data

			err := pctx.UnmarshalProviderData(&data)
			if err != nil {
				return err
			}

			// Check if the machine already exists
			vm, err := p.harvesterClient.HarvesterClient.KubevirtV1().VirtualMachines(p.namespace).Get(ctx, pctx.GetRequestID(), k8smetav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				logger.Error("failed to get the machine", zap.Error(err))

				return err
			}
			if err == nil {
				logger.Info("machine already exists", zap.String("machineName", vm.Name))

				return nil
			}

			vm = &kvv1.VirtualMachine{}
			vm.Spec.Running = pointer.To(true)

			if vm.Spec.Template == nil {
				vm.Spec.Template = &kvv1.VirtualMachineInstanceTemplateSpec{
					Spec: kvv1.VirtualMachineInstanceSpec{
						Domain: kvv1.DomainSpec{
							Resources: kvv1.ResourceRequirements{
								Requests: v1.ResourceList{},
							},
						},
					},
				}
			}

			vm.Spec.Template.Spec.Architecture = data.Architecture
			vm.Spec.Template.Spec.Domain.CPU = &kvv1.CPU{
				Cores: uint32(data.Cores),
			}

			vm.Spec.Template.Spec.Domain.Resources.Requests[v1.ResourceMemory] = *resource.NewQuantity(int64(data.Memory)*1024*1024, resource.DecimalSI)

			if vm.Spec.Template.ObjectMeta.Labels == nil {
				vm.Spec.Template.ObjectMeta.Labels = map[string]string{}
			}

			vm.Spec.Template.ObjectMeta.Labels["tag.harvesterhci.io/created-by"] = "omni-infra-provider-harvester"
			vm.Spec.Template.ObjectMeta.Labels["tag.harvesterhci.io/managed-by"] = "omni"
			vm.Spec.Template.ObjectMeta.Labels["harvesterhci.io/creator"] = "omni-infra-provider-harvester"
			vm.Spec.Template.ObjectMeta.Labels["harvesterhci.io/vmName"] = pctx.GetRequestID()

			// Set the affinity to prefer the node the machine is scheduled to
			vm.Spec.Template.Spec.Affinity = &v1.Affinity{
				PodAntiAffinity: &v1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
						{
							Weight: 100,
							PodAffinityTerm: v1.PodAffinityTerm{
								LabelSelector: &k8smetav1.LabelSelector{
									MatchExpressions: []k8smetav1.LabelSelectorRequirement{
										{
											Key:      "harvesterhci.io/creator",
											Operator: k8smetav1.LabelSelectorOpExists,
										},
									},
								},
								TopologyKey: "kubernetes.io/hostname",
							},
						},
					},
				},
			}

			// Set the firmware and disable secure boot
			vm.Spec.Template.Spec.Domain.Firmware = &kvv1.Firmware{
				UUID: types.UID(pctx.State.TypedSpec().Value.Uuid),
				Bootloader: &kvv1.Bootloader{
					EFI: &kvv1.EFI{
						SecureBoot: pointer.To(false),
					},
				},
			}

			// Set the network interface
			vm.Spec.Template.Spec.Networks = []kvv1.Network{
				{
					Name: "default",
					NetworkSource: kvv1.NetworkSource{
						Multus: &kvv1.MultusNetwork{
							NetworkName: fmt.Sprintf("%s/%s", data.NetworkNamespace, data.NetworkName),
						},
					},
				},
			}

			// Set the disks and volumes
			vm.Spec.Template.Spec.Domain.Devices.Interfaces = []kvv1.Interface{
				*kvv1.DefaultBridgeNetworkInterface(),
			}

			vm.Spec.Template.Spec.Domain.Devices.Disks = []kvv1.Disk{
				{
					Name:      "disk0",
					BootOrder: pointer.To(uint(1)),
					DiskDevice: kvv1.DiskDevice{
						Disk: &kvv1.DiskTarget{
							Bus: kvv1.DiskBusVirtio,
						},
					},
				},
				{
					Name: "cloudinitdisk",
					DiskDevice: kvv1.DiskDevice{
						Disk: &kvv1.DiskTarget{
							Bus: kvv1.DiskBusVirtio,
						},
					},
				},
			}

			vm.Spec.Template.Spec.Volumes = []kvv1.Volume{
				{
					Name: "disk0",
					VolumeSource: kvv1.VolumeSource{
						PersistentVolumeClaim: &kvv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: fmt.Sprintf("%s-%s-%s", pctx.GetRequestID(), "disk-0", pctx.State.TypedSpec().Value.Uuid[0:8]),
							},
						},
					},
				},
				{
					Name: "cloudinitdisk",
					VolumeSource: kvv1.VolumeSource{
						CloudInitNoCloud: &kvv1.CloudInitNoCloudSource{
							UserData:    "#cloud-config\npassword: talos\nchpasswd: { expire: False }\nssh_pwauth: True\n",
							NetworkData: "",
						},
					},
				},
			}

			// Set default input devices
			vm.Spec.Template.Spec.Domain.Devices.Inputs = []kvv1.Input{
				{
					Name: "tablet",
					Bus:  kvv1.InputBusUSB,
					Type: kvv1.InputTypeTablet,
				},
			}

			// Set the PVC annotation
			if vm.Spec.Template.ObjectMeta.Annotations == nil {
				vm.Spec.Template.ObjectMeta.Annotations = map[string]string{}
			}

			var pvcAnnotation PVCRequest
			pvcAnnotation.Metadata.Name = fmt.Sprintf("%s-%s-%s", pctx.GetRequestID(), "disk-0", pctx.State.TypedSpec().Value.Uuid[0:8])
			pvcAnnotation.Metadata.Annotations = map[string]string{
				"harvesterhci.io/imageId": fmt.Sprintf("%s/%s", p.namespace, pctx.State.TypedSpec().Value.VolumeId),
			}
			pvcAnnotation.Spec.AccessModes = []string{"ReadWriteOnce"}
			pvcAnnotation.Spec.Resources.Requests.Storage = fmt.Sprintf("%dGi", data.DiskSize)
			pvcAnnotation.Spec.VolumeMode = "Block"
			pvcAnnotation.Spec.StorageClassName = "longhorn-" + pctx.State.TypedSpec().Value.VolumeId

			annotation, err := pvcAnnotation.String()
			if err != nil {
				logger.Error("failed to marshal PVC annotation", zap.Error(err))

				return err
			}

			vm.Spec.Template.ObjectMeta.Annotations["harvesterhci.io/volumeClaimTemplates"] = annotation

			if vm.Name == "" {
				vm.Name = pctx.GetRequestID()
				vm.Namespace = p.namespace

				_, err := p.harvesterClient.
					HarvesterClient.KubevirtV1().
					VirtualMachines(p.namespace).
					Create(ctx, vm, k8smetav1.CreateOptions{})
				if err != nil {
					return err
				}
			} else {
				_, err := p.harvesterClient.
					HarvesterClient.
					KubevirtV1().
					VirtualMachines(p.namespace).
					Update(ctx, vm, k8smetav1.UpdateOptions{})
				if err != nil {
					return err
				}
			}

			return provision.NewRetryInterval(time.Second * 10)
		}),
	}
}

// Deprovision implements infra.Provisioner.
func (p *Provisioner) Deprovision(ctx context.Context, logger *zap.Logger, _ *resources.Machine, machineRequest *infra.MachineRequest) error {
	logger.Info("deprovisioning machine", zap.String("machineName", machineRequest.Metadata().ID()))

	_, err := p.harvesterClient.HarvesterClient.KubevirtV1().
		VirtualMachines(p.namespace).Get(ctx, machineRequest.Metadata().ID(), k8smetav1.GetOptions{})

	if errors.IsNotFound(err) {
		logger.Info("machine not found, skipping deletion")

		return nil
	}

	if err != nil {
		logger.Error("failed to get the machine", zap.Error(err))

		return provision.NewRetryInterval(time.Second * 5)
	}

	err = p.harvesterClient.HarvesterClient.KubevirtV1().
		VirtualMachines(p.namespace).Delete(ctx, machineRequest.Metadata().ID(), k8smetav1.DeleteOptions{
		PropagationPolicy: pointer.To(k8smetav1.DeletePropagationForeground),
	})

	if err != nil && !errors.IsNotFound(err) {
		return provision.NewRetryInterval(time.Second * 5)
	}

	logger.Info("machine deleted", zap.String("machineName", machineRequest.Metadata().ID()))

	return nil
}
