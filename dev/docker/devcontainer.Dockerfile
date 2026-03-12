FROM golang:1.26

# Install golangci-lint
RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.3

# Add spellcheck and jq
RUN apt-get update && apt-get install -y \
    shellcheck \
    jq
