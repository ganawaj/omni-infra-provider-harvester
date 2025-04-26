// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

// Data is the provider custom machine config.
type Data struct {
	Architecture     string `yaml:"architecture"`
	StorageClass     string `yaml:"storage_class"`
	NetworkName      string `yaml:"network_name"`
	NetworkNamespace string `yaml:"network_namespace"`
	Namespace        string `yaml:"namespace"`
	Memory           uint64 `yaml:"memory"`
	Cores            int    `yaml:"cores"`
	DiskSize         int    `yaml:"disk_size"`
}
