package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"

	"nacl-cli/box"
)

type version subcommands.Commander

func (flg *version) Name() string           { return "version" }
func (flg *version) Synopsis() string       { return "show version" }
func (flg *version) SetFlags(*flag.FlagSet) {}
func (flg *version) Usage() string          { return "version:\n\tshow version\n" }
func (flg *version) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Printf("cli: v0.0.0\nnacl: v0.0.0\n")
	return subcommands.ExitSuccess
}
func versionCommand() subcommands.Command {
	return (*version)(subcommands.DefaultCommander)
}


func main() {
	subcommands.Register(subcommands.HelpCommand(), "misc")
	subcommands.Register(versionCommand(), "misc")
	subcommands.Register(&box.Cmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
