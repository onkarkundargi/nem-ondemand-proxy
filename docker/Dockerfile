# Copyright 2016 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# -------------
# Build stage

FROM golang:1.13.8-alpine3.11 AS build-env

# Install required packages
RUN apk add --no-cache wget git make build-base protobuf protobuf-dev

WORKDIR /go/src

COPY . ./

ARG org_label_schema_version=unknown
ARG org_label_schema_vcs_url=unknown
ARG org_label_schema_vcs_ref=unknown
ARG org_label_schema_build_date=unknown
ARG org_opencord_vcs_commit_date=unknown
ARG org_opencord_vcs_dirty=unknown

# Build nem-proxy
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN go build -mod=vendor -o /build/nem-proxy cmd/nem-ondemand-proxy/main.go
# -------------
# Image creation stage

FROM alpine:3.9.4

# Set the working directory
WORKDIR /app

# Copy required files
COPY --from=build-env /build/nem-proxy /app/

# Label image
ARG org_label_schema_version=unknown
ARG org_label_schema_vcs_url=unknown
ARG org_label_schema_vcs_ref=unknown
ARG org_label_schema_build_date=unknown
ARG org_opencord_vcs_commit_date=unknown
ARG org_opencord_vcs_dirty=unknown

LABEL org.label-schema.schema-version=1.0 \
      org.label-schema.name=nem-proxy \
      org.label-schema.version=$org_label_schema_version \
      org.label-schema.vcs-url=$org_label_schema_vcs_url \
      org.label-schema.vcs-ref=$org_label_schema_vcs_ref \
      org.label-schema.build-date=$org_label_schema_build_date \
      org.opencord.vcs-commit-date=$org_opencord_vcs_commit_date \
      org.opencord.vcs-dirty=$org_opencord_vcs_dirty

CMD /app/nem-proxy
