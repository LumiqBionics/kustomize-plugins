FROM golang:1.12 as builder

COPY . /code
WORKDIR /code

ENV GOOS linux
ENV GOARCH amd64
ENV GOPATH /gopath

RUN mkdir -p /gopath && \
    go install sigs.k8s.io/kustomize/v3/cmd/kustomize && \
    go build -buildmode plugin -o SecretsFromVault.so SecretsFromVault.go


FROM debian:9.5-slim

ENV XDG_CONFIG_HOME=/xdg

COPY --from=builder /code/SecretsFromVault.so \
    /xdg/kustomize/plugin/lumiq.com/v1/secretsfromvault/SecretsFromVault.so
COPY --from=builder /gopath/bin/kustomize /kustomize

ENTRYPOINT ["/kustomize"]

WORKDIR /code

CMD ["build", "--enable_alpha_plugins"]
