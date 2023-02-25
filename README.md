# xmtpd

XMTP node implementation.

[![Test](https://github.com/xmtp/xmtpd/actions/workflows/test.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/test.yml)
[![Build](https://github.com/xmtp/xmtpd/actions/workflows/build.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/build.yml)
[![Publish](https://github.com/xmtp/xmtpd/actions/workflows/publish.yml/badge.svg)](https://github.com/xmtp/xmtpd/actions/workflows/publish.yml)

## Development

Build and install dependencies:

```sh
dev/up
```

Start a node:

```sh
dev/start
```

Run tests:

```sh
dev/test
```

Visit local [Prometheus](https://prometheus.io/) UI:
```sh
open http://localhost:9090
```

Visit local [Jaeger](https://www.jaegertracing.io/) UI:
```sh
open http://localhost:16686
```

## Resources

* OpenTelemetry intro https://www.komu.engineer/blogs/11/opentelemetry-and-go
