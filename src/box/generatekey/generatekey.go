package generatekey

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/nacl/box"
	"github.com/google/subcommands"
)

type Cmd struct {
	keyFilePath string
	publicKeyFileSuffix string
	privateKeyFileSuffix string
}

func (*Cmd) Name() string     { return "generatekey" }
func (*Cmd) Synopsis() string { return "Generate a public key and private key" }
func (*Cmd) Usage() string {
	return `box generatekey -name key:
	Generate a public key and private key to key.pub and key respectively.

	see https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc#GenerateKey

`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
  f.StringVar(&c.keyFilePath, "output", "sender", "path to the generated keys")
  f.StringVar(&c.publicKeyFileSuffix, "suffix", ".pub", "suffix of public key file")
  f.StringVar(&c.privateKeyFileSuffix, "private-suffix", ".key", "suffix of private key file")
}

func (c *Cmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)

	if err != nil {
		panic(err)
	}

	if dir, _ := filepath.Split(c.keyFilePath); dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}
	if err := ioutil.WriteFile(c.keyFilePath + c.publicKeyFileSuffix, publicKey[:], 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return subcommands.ExitFailure
	}
	if err := ioutil.WriteFile(c.keyFilePath + c.privateKeyFileSuffix, privateKey[:], 0640); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
