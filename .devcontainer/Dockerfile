FROM debian:stable-slim

ARG GO_VERSION=1.23.0
ARG PYTHON_VERSION=3.12.4

ENV GOPATH=/go
ENV GOBIN=${GOPATH}/bin
ENV CGO_ENABLED=1
ENV PATH=${GOBIN}:${GOPATH}/bin:/usr/local/go/bin:${PATH}

RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    python3 \
    python3-pip \
    python3-venv \
    python3-ipykernel

ENV VENV_PATH=/opt/venv
RUN python3 -m venv $VENV_PATH
ENV PATH="${VENV_PATH}/bin:${PATH}"

RUN pip install --no-cache-dir ipykernel pyarrow pandas

RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz" | tar -C /usr/local -xz

RUN go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

RUN rm -rf /var/lib/apt/lists/*

# TODO will probably need a go get step in here to install go project dependencies of the project

RUN useradd -ms /bin/bash vscode
RUN chown -R vscode:vscode /go
USER vscode

WORKDIR /workspaces

CMD ["bash"]