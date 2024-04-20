//go:build tools
// +build tools

package internal

import (
	_ "github.com/bokwoon95/wgo"
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "github.com/golangci/golangci-lint/pkg/commands"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/stephenafamo/bob/gen/bobgen-sqlite"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
