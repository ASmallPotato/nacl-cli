```sh
./bin/nacl box generatekey -output=keys/sender
./bin/nacl box generatekey -output=keys/recipient

echo 'hello world' | ./bin/nacl box seal -private-key=keys/sender.key -public-key=keys/recipient.pub | tee /tmp/encrypted_message
./bin/nacl box open -private-key=keys/recipient.key -public-key=keys/sender.pub < /tmp/encrypted_message
```
