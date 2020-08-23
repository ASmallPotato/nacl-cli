cli for [NaCl](https://nacl.cr.yp.to)
====

This is a thin wrapper around [go's nacl](https://pkg.go.dev/golang.org/x/crypto/nacl), as a simple command line binding to NaCl.

[saltpack](https://github.com/keybase/saltpack) should satisfy most cryptography needs and should be used instead of this program.
This is most useful for generating keys manually to communicate with other applications.

current only [box](https://pkg.go.dev/golang.org/x/crypto/nacl/box?tab=doc) is supported (and only GenerateKey, Open and Seal are implemented).
