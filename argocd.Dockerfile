ARG ARGO_VERSION=v1.1.0

FROM golang:1.12 as builder

COPY . /code
WORKDIR /code

ENV GOOS linux
ENV GOARCH amd64
ENV GOPATH /gopath

RUN mkdir -p /gopath && \
    go install sigs.k8s.io/kustomize/v3/cmd/kustomize && \
    go build -buildmode plugin -o SecretsFromVault.so SecretsFromVault.go


FROM argoproj/argocd:$ARGO_VERSION

COPY --from=builder /code/SecretsFromVault.so \
    /home/argocd/.config/kustomize/plugin/lumiq.com/v1/secretsfromvault/SecretsFromVault.so
COPY --from=builder /gopath/bin/kustomize /usr/local/bin/kustomize3
COPY scripts/kustomize /usr/local/bin/kustomize