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
	return `box generatekey -output=keys/sender:
	Generate a public key and private key to key.pub and key respectively.

	to output to stdout use extra file descriptor:
		nacl box generatekey -output=/ -suffix >(base64) -private-suffix >(base64)

	see https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc#GenerateKey

`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
  f.StringVar(&c.keyFilePath, "output", "sender", "path to the generated keys")
  f.StringVar(&c.publicKeyFileSuffix, "suffix", ".pub", "suffix of public key file")
  f.StringVar(&c.privateKeyFileSuffix, "private-suffix", ".key", "suffix of private key file")
}

func (c *Cmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var filePath string
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)

	if err != nil {
		panic(err)
	}

	if dir, _ := filepath.Split(c.keyFilePath); dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}
	filePath = c.keyFilePath + c.publicKeyFileSuffix
	if err := ioutil.WriteFile(filePath, publicKey[:], 0644); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return subcommands.ExitFailure
	}
	filePath = c.keyFilePath + c.privateKeyFileSuffix
	if err := ioutil.WriteFile(filePath, privateKey[:], 0640); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
