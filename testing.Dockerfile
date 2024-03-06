FROM golang:1.22
ARG TARGETOS
ARG TARGETARCH

RUN apt update && \
apt upgrade -y && \
apt install dumb-init kubernetes-client -y && \
apt install -y\
  btrfs-progs \
  crun \
  git \
  golang-go \
  go-md2man \
  iptables \
  libassuan-dev \
  libbtrfs-dev \
  libc6-dev \
  libdevmapper-dev \
  libglib2.0-dev \
  libgpgme-dev \
  libgpg-error-dev \
  libprotobuf-dev \
  libprotobuf-c-dev \
  libseccomp-dev \
  libselinux1-dev \
  libsystemd-dev \
  netavark \
  pkg-config \
  uidmap

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest && setup-envtest use
RUN mkdir -p /root/.kube && touch /root/.kube/config

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download