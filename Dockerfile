FROM golang:1.12 as builder

ARG KUSTOMIZE_VERSION=v3.1.0
ARG KUBE_VERSION=v1.14.0
COPY . /code
WORKDIR /code

ENV GOOS linux
ENV GOARCH amd64
ENV GOPATH /gopath

RUN mkdir -p /gopath sigs.k8s.io \
    && git clone https://github.com/kubernetes-sigs/kustomize.git sigs.k8s.io/kustomize \
    && (cd sigs.k8s.io/kustomize; git checkout ${KUSTOMIZE_VERSION}) \
    && cp SecretsFromVault.go sigs.k8s.io/kustomize/plugin/ \
    && cd sigs.k8s.io/kustomize \
    && git apply ../../kustomize.patch \
    && git apply ../../kustomize-enable.patch \
    && go build -o $GOPATH/bin/kustomize cmd/kustomize/main.go \
    && go build -buildmode plugin -o /SecretsFromVault.so plugin/SecretsFromVault.go \
    && wget https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl \
        -O /kubectl \
    && chmod a+x /kubectl

FROM debian:9.5-slim

ENV XDG_CONFIG_HOME=/xdg

COPY --from=builder /SecretsFromVault.so \
    /xdg/kustomize/plugin/lumiq.com/v1/secretsfromvault/SecretsFromVault.so
COPY --from=builder /gopath/bin/kustomize /kustomize
COPY --from=builder /kubectl /kubectl

ENTRYPOINT ["/kustomize"]

WORKDIR /code

CMD ["build"]
