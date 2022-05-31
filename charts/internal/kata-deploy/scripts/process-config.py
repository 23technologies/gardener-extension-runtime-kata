# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

import toml

toml_dict = toml.load("/etc/containerd/config.toml", _dict=dict)

toml_dict['plugins']['io.containerd.snapshotter.v1.devmapper']['base_image_size'] = '40GB'
toml_dict['plugins']['io.containerd.snapshotter.v1.devmapper']['root_path'] = '/var/lib/containerd/devmapper'
toml_dict['plugins']['io.containerd.snapshotter.v1.devmapper']['pool_name'] = 'devpool'

toml_dict['plugins']['io.containerd.grpc.v1.cri']['containerd']['snapshotter'] = 'devmapper'

with open('/etc/containerd/config.toml', 'w') as f:
    toml.dump(toml_dict, f)
