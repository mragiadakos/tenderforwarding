Tenderforwarding

Tenderforwarding has only one purpose, for the inflator to send terndermoney's coins from the taxes to other people.
It will do that by encrypting the private key of the tendermoney's coin based on the public key of the receiver using ECDH.
For that reason, it will have only one action
- FORWARD:
  It will send money, but only the redistributor can send money to other people with metadata that will be used for reference


Delivery
Request:
{
    Redistributor: Public key hex
    Signature: hex
    Coins: []uuid
    EncryptedPrivateKeys: a hex of the encrypted json map[uuid]private_key_hex
    Metadata: map[string]string
}
Response:
  The request will fail if:
  - The redistributor is not in the list
  - The signature is not correct
  - The coins are empty
  - The encrypted data is empty
  - The metadata is empty