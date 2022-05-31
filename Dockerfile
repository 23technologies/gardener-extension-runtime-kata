# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

############# builder
FROM eu.gcr.io/gardener-project/3rd/golang:1.17.8 AS builder

ENV BINARY_PATH=/go/bin
WORKDIR /go/src/github.com/23technologies/gardener-extension-runtime-kata

COPY . .
RUN make install

############# base
FROM eu.gcr.io/gardener-project/3rd/alpine:3.15 as base

############# gardener-extension-runtime-kata
FROM base AS gardener-extension-runtime-kata
LABEL org.opencontainers.image.source="https://github.com/23technologies/gardener-extension-runtime-kata"

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-runtime-kata /gardener-extension-runtime-kata
ENTRYPOINT ["/gardener-extension-runtime-kata"]
