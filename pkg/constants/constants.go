// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package constants defines constants used
package constants

const (
	// ExtensionType is the name of the extension type.
	ExtensionType = "runtime-kata"
	// ServiceName is the name of the service.
	ServiceName = ExtensionType

	extensionServiceName = "extension-" + ServiceName

	ManagedResourceNameKataResources = extensionServiceName + "-resources" 
	ManagedResourceNameKataDeploy = extensionServiceName + "-deploy" 
	ManagedResourceNameKataCleanup = extensionServiceName + "-cleanup" 

	ReleaseNameKataResources = "kata-resources"
	ReleaseNameKataDeploy = "kata-deploy"
	ReleaseNameKataCleanup = "kata-cleanup"

	ChartPathKataResources = "charts/internal/kata-resources"
	ChartPathKataDeploy = "charts/internal/kata-deploy"
	ChartPathKataCleanup = "charts/internal/kata-cleanup"

	KataSetupFcScriptPath = "/var/lib/setup-kata-fc.sh"
	KataSetupFcUnitName = "devmapper-reload.service"
)
