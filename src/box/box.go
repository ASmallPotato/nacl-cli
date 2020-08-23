package box

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"nacl-cli/box/generatekey"
	"nacl-cli/box/open"
	"nacl-cli/box/seal"
)

type Cmd struct {
}

func (*Cmd) Name() string     { return "box" }
func (*Cmd) Synopsis() string { return "authenticates and encrypts small messages using public-key cryptography" }
func (*Cmd) Usage() string {
	return `box:
	see https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc

`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
}

func (c *Cmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	commander := subcommands.NewCommander(f, "box")

	commander.Register(&generatekey.Cmd{}, "")
	commander.Register(&seal.Cmd{}, "")
	commander.Register(&open.Cmd{}, "")

	return commander.Execute(ctx, args...)
}
