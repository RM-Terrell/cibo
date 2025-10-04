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
    # python, pip, venv, ipykernel, pyarrow, and pandas are all installed for Data Wrangler to work
    python3 \
    python3-pip \
    python3-venv

ENV VENV_PATH=/opt/venv
RUN python3 -m venv $VENV_PATH
ENV PATH="${VENV_PATH}/bin:${PATH}"

RUN pip install --no-cache-dir ipykernel pyarrow pandas

RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz" | tar -C /usr/local -xz

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

RUN go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

RUN rm -rf /var/lib/apt/lists/*

RUN useradd -ms /bin/bash vscode
# need to provide perms on /go or else vscode go extensions wont work
RUN chown -R vscode:vscode /go
# need to provide perms on /venv or else vscode Python extensions wont be able to read the venv
RUN chown -R vscode:vscode /opt/venv
USER vscode

WORKDIR /workspaces

CMD ["bash"]