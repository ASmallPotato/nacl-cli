package open

import (
	"context"
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
	encryptedMessageFilePath string
	// nonceFilePath string
}

func (*Cmd) Name() string     { return "open" }
func (*Cmd) Synopsis() string { return "open message sign by sender's public key with private key" }
func (*Cmd) Usage() string {
	return `box open -private-key=/tmp/key -public-key=/tmp/key.pub < /tmp/message:
	Authenticates and decrypts input produced by seal.

	see https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc#Open

`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.publicKeyFilePath, "public-key", "", "path to public key")
	f.StringVar(&c.privateKeyFilePath, "private-key", "", "path to private key")
	f.StringVar(&c.encryptedMessageFilePath, "file", "-", "path to input file, use - for stdin")
	// f.StringVar(&c.nonceFilePath, "nonce", "", "path of file to read nonce from if provided")
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
	var encryptedMessage []byte
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

	if reader, err = resolveFileFromPath(c.encryptedMessageFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read message: %s\n", err)
		return subcommands.ExitUsageError
	}
	if encryptedMessage, err = ioutil.ReadAll(reader); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read message: %s\n", err)
		return subcommands.ExitUsageError
	}
	if len(encryptedMessage) < 24 {
		msg := (
			"Unable to decrypt message, " +
			"open is only able to decrpty message generated from seal. " +
			"Namely message must be prefixed with 24 byte of nonce." +
		"\n")
		fmt.Fprintf(os.Stderr, msg)
		return subcommands.ExitUsageError
	}

	copy(nonce[:], encryptedMessage[:24])

	decrypted, ok := box.Open(
		nil,
		encryptedMessage[24:],
		&nonce,
		&publicKey,
		&privateKey,
	)
	if !ok {
		fmt.Fprintf(os.Stderr, "decryption error\n")
		return subcommands.ExitFailure
	}

	fmt.Printf("%s", decrypted)

	return subcommands.ExitSuccess
}
