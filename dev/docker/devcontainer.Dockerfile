FROM mcr.microsoft.com/vscode/devcontainers/go:1.25

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2

# Add spellcheck and jq
RUN apt-get update && apt-get install -y \
    shellcheck \
    jq
