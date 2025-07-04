FROM debian:stable-slim

ARG GO_VERSION=1.23.0

ENV GOPATH=/go
ENV GOBIN=${GOPATH}/bin
ENV CGO_ENABLED=1
ENV PATH=${GOBIN}:${GOPATH}/bin:/usr/local/go/bin:${PATH}

RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential

RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz" | tar -C /usr/local -xz

RUN go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

RUN useradd -ms /bin/bash vscode
USER vscode

WORKDIR /workspaces

CMD ["bash"]