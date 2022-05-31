#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -ex

yum -y install python3 bc
pip3 install toml

python3 /scripts-dir/process-config.py

#systemctl restart containerd
