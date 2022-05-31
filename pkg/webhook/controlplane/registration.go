// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/extensions/pkg/webhook/controlplane"
	"github.com/gardener/gardener/extensions/pkg/webhook/controlplane/genericmutator"


	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	oscutils "github.com/gardener/gardener/pkg/operation/botanist/component/extensions/operatingsystemconfig/utils"
	"github.com/gardener/gardener/pkg/operation/botanist/component/extensions/operatingsystemconfig/original/components/kubelet"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var logger = log.Log.WithName("runtime-kata")

// AddToManager creates a webhook and adds it to the manager.
//
// PARAMETERS
// mgr manager.Manager Webhook control plane controller manager instance
func AddToManager(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Adding webhook to manager")
	fciCodec := oscutils.NewFileContentInlineCodec()
	
	logger := logger.WithValues("kind", controlplane.KindSeed)

	mutator := genericmutator.NewMutator(
		NewEnsurer(logger),
		oscutils.NewUnitSerializer(),
		kubelet.NewConfigCodec(fciCodec),
		fciCodec,
		logger,
	)

	Types :=    []extensionswebhook.Type{
			{ Obj: &extensionsv1alpha1.OperatingSystemConfig{} },
		}

	// Create handler
	handler, err := extensionswebhook.NewBuilder(mgr, logger).WithMutator(mutator, Types...).Build()
	if err != nil {
		return nil, err
	}

	// Create webhook
	logger.Info("Creating webhook", "name", "runtime-kata")

	return &extensionswebhook.Webhook{
		Name:    "controlplaneexposure",
		Kind:    controlplane.KindSeed,
		Types:    Types,
		Target:   extensionswebhook.TargetSeed,
		Path:     "controlplaneexposure",
		Webhook:  &admission.Webhook{Handler: handler},
	}, nil

	return controlplane.New(mgr, controlplane.Args{
		Kind:    controlplane.KindSeed,
		Provider:  "rutime-kata",
		Types:    []extensionswebhook.Type{
			{ Obj: &extensionsv1alpha1.OperatingSystemConfig{} },
		},
		Mutator: genericmutator.NewMutator(
			NewEnsurer(logger),
			oscutils.NewUnitSerializer(),
			kubelet.NewConfigCodec(fciCodec),
			fciCodec,
			logger,
		),
	})
}
