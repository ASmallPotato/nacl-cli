package seal

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/nacl/box"
	"github.com/google/subcommands"
)

type Cmd struct {
	publicKeyFilePath string
	privateKeyFilePath string
	messageFilePath string
	nonceFilePath string
}

func (*Cmd) Name() string     { return "seal" }
func (*Cmd) Synopsis() string { return "seal message with private key for public key peer's viewing" }
func (*Cmd) Usage() string {
	return `box seal -public-key=/tmp/recipient.pub -private-key=/tmp/sender < /tmp/message:
	appends an encrypted input, prepended with nonce. The nonce must be unique for each distinct message for a given pair of keys.

	see https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc#Seal

`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
  f.StringVar(&c.publicKeyFilePath, "public-key", "", "path to public key")
  f.StringVar(&c.privateKeyFilePath, "private-key", "", "path to private key")
  f.StringVar(&c.messageFilePath, "file", "-", "path to input file, use - for stdin")
  f.StringVar(&c.nonceFilePath, "nonce", "", "path of file to read nonce from if provided")
}

func resolveFileFromPath(path string) (io.Reader, error) {
	switch path {
	case "-":
		return os.Stdin, nil
	case "":
		return nil, fmt.Errorf("path not provided\n")
	default:
		return os.Open(path)
	}
}

func (c *Cmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var publicKey [32]byte
	var privateKey [32]byte
	var message []byte
	var nonce [24]byte

	var reader io.Reader
	var err error

	if reader, err = resolveFileFromPath(c.publicKeyFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read public key: %s\n", err)
		return subcommands.ExitUsageError
	}
	io.ReadFull(reader, publicKey[:])

	if reader, err = resolveFileFromPath(c.privateKeyFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read private key: %s\n", err)
		return subcommands.ExitUsageError
	}
	io.ReadFull(reader, privateKey[:])

	if reader, err = resolveFileFromPath(c.messageFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read message: %s\n", err)
		return subcommands.ExitUsageError
	}
	if message, err = ioutil.ReadAll(reader); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read message: %s\n", err)
		return subcommands.ExitUsageError
	}

	if c.nonceFilePath != "" {
		if reader, err = resolveFileFromPath(c.nonceFilePath); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read nonce: %s\n", err)
			return subcommands.ExitUsageError
		}
	} else {
		reader = rand.Reader
	}
	io.ReadFull(reader, nonce[:])


	encrypted := box.Seal(
		nonce[:],
		message,
		&nonce,
		&publicKey,
		&privateKey,
	)

	fmt.Printf("%s", encrypted)

	return subcommands.ExitSuccess
}
