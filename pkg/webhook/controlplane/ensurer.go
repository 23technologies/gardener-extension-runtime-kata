// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"context"

	"github.com/23technologies/gardener-extension-runtime-kata/pkg/constants"
	gcontext "github.com/gardener/gardener/extensions/pkg/webhook/context"
	"github.com/gardener/gardener/extensions/pkg/webhook/controlplane/genericmutator"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/utils/pointer"
)

// NewEnsurer creates a new controlplane ensurer.
 func NewEnsurer(logger logr.Logger) genericmutator.Ensurer {
	 return &ensurer{
		 logger: logger.WithName("azure-controlplane-ensurer"),
	 }
 }

 type ensurer struct {
	 genericmutator.NoopEnsurer
	 client client.Client
	 logger logr.Logger
 }

 // InjectClient injects the given client into the ensurer.
 func (e *ensurer) InjectClient(client client.Client) error {
	 e.client = client
	 return nil
 }

func (e *ensurer) EnsureAdditionalUnits(ctx context.Context, gctx gcontext.GardenContext, new, old *[]extensionsv1alpha1.Unit) error {

	// check if this path is already in the current OperatingSystemConfig
	// if so, return, as otherwise the shoot reconcilation will fail due to duplicate entries in the osc
	for _, val := range(*new) {
		if val.Name == constants.KataSetupFcUnitName {
			return nil
		}
	}

	unitContent := `[Unit]
Description=Devmapper reload script

[Service]
ExecStart=` + constants.KataSetupFcScriptPath + `

[Install]
WantedBy=multi-user.target
`
	
	unit := extensionsv1alpha1.Unit{
		Name:   constants.KataSetupFcUnitName, 
		Command: nil, 
		Enable: pointer.Bool(true),
		Content: pointer.String(unitContent), 
		DropIns: nil,
	}

	*new = append(*new, unit)

	return nil
}

// EnsureAdditionalFiles ensures that additional required system files are added.
func (e *ensurer) EnsureAdditionalFiles(ctx context.Context, gctx gcontext.GardenContext, new, old *[]extensionsv1alpha1.File) error {

	// set the scriptPath for our script to be saved on the node
	scriptPath := constants.KataSetupFcScriptPath 

	// check if this path is already in the current OperatingSystemConfig
	// if so, return, as otherwise the shoot reconcilation will fail due to duplicate entries in the osc
	for _, val := range(*new) {
		if val.Path == scriptPath {
			return nil
		}
	}

	// define the script. You can find some explanation here
	// https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-kata-containers-with-firecracker.md
	script := `#!/bin/bash
set -ex

DATA_DIR=/var/lib/containerd/devmapper
POOL_NAME=devpool

mkdir -p ${DATA_DIR}

# Create data file
touch ${DATA_DIR}/data
truncate -s 100G ${DATA_DIR}/data

# Create metadata file
touch ${DATA_DIR}/meta
truncate -s 10G ${DATA_DIR}/meta

# Allocate loop devices
DATA_DEV=$(sudo losetup --find --show ${DATA_DIR}/data)
META_DEV=$(sudo losetup --find --show ${DATA_DIR}/meta)

# Define thin-pool parameters.
# See https://www.kernel.org/doc/Documentation/device-mapper/thin-provisioning.txt for details.
SECTOR_SIZE=512
DATA_SIZE=$(sudo blockdev --getsize64 -q ${DATA_DEV})
LENGTH_IN_SECTORS=$(bc <<< ${DATA_SIZE}/${SECTOR_SIZE})
DATA_BLOCK_SIZE=128
LOW_WATER_MARK=32768

# Create a thin-pool device
dmsetup create ${POOL_NAME} \
    --table "0 ${LENGTH_IN_SECTORS} thin-pool ${META_DEV} ${DATA_DEV} ${DATA_BLOCK_SIZE} ${LOW_WATER_MARK}"
`

  scriptFile := extensionsv1alpha1.File{
  	Path:        scriptPath,
  	Permissions: pointer.Int32Ptr(0755),
  	Content:     extensionsv1alpha1.FileContent{
  		SecretRef:         nil,
  		Inline:            &extensionsv1alpha1.FileContentInline{
  			Encoding: "",
  			Data:     script,
  		},
  		TransmitUnencoded: nil,
  	},
  }

	*new = append(*new, scriptFile)
	
	return nil
}

