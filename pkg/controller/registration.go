// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package controller 

import (
	"time"

	"github.com/gardener/gardener/extensions/pkg/controller/containerruntime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	resourcesv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
)

// DefaultAddOptions contains configuration for the mwe controller 
var DefaultAddOptions = AddOptions{}

// AddOptions are options to apply when adding the mwe controller to the manager.
type AddOptions struct {
	// ControllerOptions contains options for the controller.
	ControllerOptions controller.Options
	// IgnoreOperationAnnotation specifies whether to ignore the operation annotation or not.
	IgnoreOperationAnnotation bool
}

// AddToManager adds a mwe Lifecycle controller to the given Controller Manager.
func AddToManager(mgr manager.Manager) error {

	scheme := mgr.GetScheme()
	if err := resourcesv1alpha1.AddToScheme(scheme); err != nil {
		return err
	}
	return containerruntime.Add(mgr, containerruntime.AddArgs{
		Actuator:          NewActuator(),
		ControllerOptions: DefaultAddOptions.ControllerOptions,
		FinalizerSuffix:   "extension-runtime-kata",
		Resync:            60 * time.Minute,
		Predicates:        containerruntime.DefaultPredicates(DefaultAddOptions.IgnoreOperationAnnotation),
		Type:              "runtime-kata",
	})
}
