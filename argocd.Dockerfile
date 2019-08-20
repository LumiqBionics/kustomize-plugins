ARG ARGO_VERSION=v1.1.0

FROM golang:1.12 as builder

ARG KUSTOMIZE_VERSION=v3.1.0
COPY . /code
WORKDIR /code

ENV GOOS linux
ENV GOARCH amd64
ENV GOPATH /gopath

RUN mkdir -p /gopath sigs.k8s.io && \
    git clone https://github.com/kubernetes-sigs/kustomize.git sigs.k8s.io/kustomize && \
    (cd sigs.k8s.io/kustomize; git checkout ${KUSTOMIZE_VERSION}) && \
    cp SecretsFromVault.go sigs.k8s.io/kustomize/plugin/ && \
    cd sigs.k8s.io/kustomize && \
    git apply ../../kustomize.patch && \
    git apply ../../kustomize-enable.patch && \
    go build -o $GOPATH/bin/kustomize cmd/kustomize/main.go && \
    go build -buildmode plugin -o /SecretsFromVault.so plugin/SecretsFromVault.go


FROM argoproj/argocd:$ARGO_VERSION

ENV XDG_CONFIG_HOME=/xdg

COPY --from=builder /SecretsFromVault.so \
    /xdg/kustomize/plugin/lumiq.com/v1/secretsfromvault/SecretsFromVault.so
COPY --from=builder /gopath/bin/kustomize /usr/local/bin/kustomize
