# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

package cmd

import (
	katawebhook "github.com/23technologies/gardener-extension-runtime-kata/pkg/webhook/controlplane"
	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
	webhook "github.com/gardener/gardener/extensions/pkg/webhook/controlplane"
)

// webhookSwitchOptions are the webhookcmd.SwitchOptions for the provider webhooks.
func webhookSwitchOptions() *webhookcmd.SwitchOptions {
	return webhookcmd.NewSwitchOptions(
		webhookcmd.Switch(webhook.WebhookName, katawebhook.AddToManager),
	)
}
