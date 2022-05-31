// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	_ "embed"
	"time"

	"github.com/23technologies/gardener-extension-runtime-kata/pkg/constants"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/containerruntime"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	resourcesv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	managedresources "github.com/gardener/gardener/pkg/utils/managedresources"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator() containerruntime.Actuator {
	return &actuator{
		logger:               log.Log.WithName("extension-runtime-kata"),
		chartRendererFactory: extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
	}
}

type actuator struct {
	logger               logr.Logger // logger
	chartRendererFactory extensionscontroller.ChartRendererFactory
	client               client.Client
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, cr *extensionsv1alpha1.ContainerRuntime, cluster *extensionscontroller.Cluster) error {

	// Before we do anything here, we check whether we meet the requirements for the extension to run. There are:
	// - nested virtualization needs to be enabled (hard requirement)
	// - ubuntu is used as the machine image (soft requirement, the extension should handle various machine images in the long run)
	err := a.checkClusterRequirements(cluster)
	if err != nil {
		return err
	}
	
	kubernetesVersion := cluster.Shoot.Spec.Kubernetes.Version

	err = a.createMrFromChart(ctx, constants.ChartPathKataResources , constants.ReleaseNameKataResources, map[string]interface{}{}, constants.ManagedResourceNameKataResources, cr.Namespace, kubernetesVersion)
	if err != nil {
		return err
	}

	kataChartValues := map[string]interface{}{
		"prepareKataFc": true,
	}
	err = a.createMrFromChart(ctx, constants.ChartPathKataDeploy, constants.ReleaseNameKataDeploy, kataChartValues, constants.ManagedResourceNameKataDeploy, cr.Namespace, kubernetesVersion)
	if err != nil {
		return err
	}

	return nil
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, cr *extensionsv1alpha1.ContainerRuntime, cluster *extensionscontroller.Cluster) error {

	// Before we do anything here, we check whether we meet the requirements for the extension to run. There are:
	// - nested virtualization needs to be enabled (hard requirement)
	// - ubuntu is used as the machine image (soft requirement, the extension should handle various machine images in the long run)
	err := a.checkClusterRequirements(cluster)
	if err != nil {
		return err
	}

	kubernetesVersion := cluster.Shoot.Spec.Kubernetes.Version

	err = a.createMrFromChart(ctx, constants.ChartPathKataCleanup, constants.ReleaseNameKataCleanup, map[string]interface{}{}, constants.ManagedResourceNameKataCleanup, cr.Namespace, kubernetesVersion)
	if err != nil {
		return err
	}

	twoMinutes := 2 * time.Minute
	timeoutShootCtx, cancelShootCtx := context.WithTimeout(ctx, twoMinutes)
	defer cancelShootCtx()

	managedresources.DeleteForShoot(ctx, a.client, cr.Namespace, constants.ManagedResourceNameKataCleanup)
	managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, cr.Namespace, constants.ManagedResourceNameKataCleanup)

	managedresources.DeleteForShoot(ctx, a.client, cr.Namespace, constants.ManagedResourceNameKataResources)
	managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, cr.Namespace, constants.ManagedResourceNameKataResources)

	return nil
}

// Restore the CR resource.
func (a *actuator) Restore(ctx context.Context, cr *extensionsv1alpha1.ContainerRuntime, cluster *extensionscontroller.Cluster) error {
	return a.Reconcile(ctx, cr, cluster)
}

// Migrate the CR resource.
func (a *actuator) Migrate(ctx context.Context, cr *extensionsv1alpha1.ContainerRuntime, cluster *extensionscontroller.Cluster) error {
	return a.Delete(ctx, cr, cluster)
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}


// mrFromChart ...
func (a *actuator) createMrFromChart(ctx context.Context, chartPath string, releaseName string, chartValues map[string]interface{}, mrName string, mrNamespace string, kubernetesVersion string) (error) {

	chartRenderer, err := a.chartRendererFactory.NewChartRendererForShoot(kubernetesVersion)
	if err != nil {
		return err
	}

	mr := &resourcesv1alpha1.ManagedResource{}
	if err := a.client.Get(ctx, kutil.Key(mrNamespace, mrName), mr); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		
		release, err := chartRenderer.Render(chartPath, releaseName, metav1.NamespaceSystem, chartValues)
		if err != nil {
			return err
		}

		err = managedresources.CreateForShoot(ctx , a.client , mrNamespace, mrName, false, map[string][]byte{"config.yaml": release.Manifest()})
		if err != nil {
			return err
		}

		twoMinutes := 2 * time.Minute
		timeoutShootCtx, cancelShootCtx := context.WithTimeout(ctx, twoMinutes)
		defer cancelShootCtx()

		managedresources.WaitUntilHealthy(timeoutShootCtx, a.client, mrNamespace, mrName)
	}
	return nil
}


// checkClusterRequirements ...
func (a *actuator) checkClusterRequirements(cl *extensionscontroller.Cluster) error  {
  	// check for the shoot running on azure first (this is the only tested provider at the moment)
	var err error
	if cl.Shoot.Spec.Provider.Type != "azure" {
		a.logger.Error(err,"Your Shoot needs to be running on a CSP with nested virtualization support", "Provider", cl.Shoot.Spec.Provider.Type)
		return err
	}
	
	// check for ubuntu being the machine image (this is the only tested machine image at the moment)
	for _, worker := range(cl.Shoot.Spec.Provider.Workers) {
		if worker.Machine.Image.Name != "ubuntu" {
		a.logger.Error(err,"Your Shoot needs to use Ubuntu as machine image in all worker pools for this extension to work properly", "Machine Image", worker.Machine.Image)
		}
	}

  return nil
}
