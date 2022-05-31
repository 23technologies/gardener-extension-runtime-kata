// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package cmd provides Kubernetes controller configuration structures used for command execution
package cmd

import (
	"context"
	"fmt"
	"os"

	controller "github.com/23technologies/gardener-extension-runtime-kata/pkg/controller"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	extensioncontrollercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/util"

	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"

	"github.com/spf13/cobra"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/component-base/version/verflag"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	restOpts    := &extensioncontrollercmd.RESTOptions{}
	reconcileOptions := &extensioncontrollercmd.ReconcilerOptions{}

	mgrOpts := &extensioncontrollercmd.ManagerOptions{
		LeaderElection:          true,
		LeaderElectionID:        extensioncontrollercmd.LeaderElectionNameID("test"),
		LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
		WebhookServerPort:       443,
	}
 
	// options for the webhook server
	webhookServerOptions := &webhookcmd.ServerOptions{
		Namespace: os.Getenv("WEBHOOK_CONFIG_NAMESPACE"),
	}

	controllerSwitches := extensioncontrollercmd.NewSwitchOptions(
			extensioncontrollercmd.Switch("runtime-kata", controller.AddToManager))
	controllerOptions := &extensioncontrollercmd.ControllerOptions{
			MaxConcurrentReconciles: 5,
		}
	
	webhookSwitches    := webhookSwitchOptions()
	webhookOptions     := webhookcmd.NewAddToManagerOptions("runtime-kata", webhookServerOptions, webhookSwitches)

	aggOption := extensioncontrollercmd.NewOptionAggregator(
	  restOpts,
		reconcileOptions,
		mgrOpts,
		webhookOptions,
		controllerOptions,
		controllerSwitches,
	)


	cmd := &cobra.Command{
		Use:           "gardener-extension-runtime-kata",
		Short:         "Controller allowing for the use of katacontainers in Shoot clusters",
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			if err := aggOption.Complete(); err != nil {
				return fmt.Errorf("error completing options: %s", err)
			}
			cmd.SilenceUsage = true

			util.ApplyClientConnectionConfigurationToRESTConfig(&componentbaseconfig.ClientConnectionConfiguration{
				QPS:   100.0,
				Burst: 130,
			}, restOpts.Completed().Config)

			mgrOpts := mgrOpts.Completed().Options()

			// do not enable a metrics server for the quick start
			mgrOpts.MetricsBindAddress = "0"

			mgr, err := manager.New(restOpts.Completed().Config, mgrOpts)
			if err != nil {
				return fmt.Errorf("could not instantiate controller-manager: %s", err)
			}

			if err := extensionscontroller.AddToScheme(mgr.GetScheme()); err != nil {
				return fmt.Errorf("could not update manager scheme: %s", err)
			}

			if _, _, err := webhookOptions.Completed().AddToManager(ctx, mgr); err != nil {
				return fmt.Errorf("Could not add webhooks to manager: %w", err)
			}

			if err := controllerSwitches.Completed().AddToManager(mgr); err != nil {
				return fmt.Errorf("could not add controllers to manager: %s", err)
			}
			
			if err := mgr.Start(ctx); err != nil {
				return fmt.Errorf("error running manager: %s", err)
			}

			return nil

		},
	}

	verflag.AddFlags(cmd.Flags())
	aggOption.AddFlags(cmd.Flags())

	return cmd
}

