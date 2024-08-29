//go:build tools
// +build tools

// I haven't figured out how to make VSCode not mark this file as a problem
// but it passes the lint check and builds fine.

package tools

import (
	_ "github.com/ethereum/go-ethereum/cmd/abigen" //nolint:all
	_ "github.com/segmentio/golines"               //nolint:all
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"          //nolint:all
	_ "github.com/vektra/mockery/v2"               //nolint:all
)
