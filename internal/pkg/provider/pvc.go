package provider

import "encoding/json"

// PVCRequest is the request to create a PVC in Harvester.
type PVCRequest struct {
	Metadata PVCMetadata `json:"metadata"`
	Spec     PVCSpec     `json:"spec"`
}

// PVCMetadata is the metadata for the PVC.
type PVCMetadata struct {
	Annotations map[string]string `json:"annotations"`
	Name        string            `json:"name"`
}

// PVCSpec is the spec for the PVC.
type PVCSpec struct {
	Resources        PVCResources `json:"resources"`
	VolumeMode       string       `json:"volumeMode"`
	StorageClassName string       `json:"storageClassName"`
	AccessModes      []string     `json:"accessModes"`
}

// PVCResources is the resources for the PVC.
type PVCResources struct {
	Requests PVCRequests `json:"requests"`
}

// PVCRequests is the requests for the PVC.
type PVCRequests struct {
	Storage string `json:"storage"`
}

func (p *PVCRequest) String() (string, error) {
	out, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
