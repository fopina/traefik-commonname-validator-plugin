//go:build tools

package tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0
//go:generate go install github.com/traefik/yaegi/cmd/yaegi@v0.14.2
