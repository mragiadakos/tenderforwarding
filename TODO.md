Tenderforwarding

Tenderforwarding has only one purpose, for the inflator to send terndermoney's coins from the taxes to other people.
It will do that by encrypting the private key of the tendermoney's coin based on the public key of the receiver using ECDH.
For that reason, it will have only one action
- FORWARD:
  It will send money, but only the redistributor can send money to other people with metadata that will be used for reference
- RECEIVED
  It will inform that the money have been received


Delivery
Request Forward:
{
    Date: time
    Signature: hex
    Action: action
    Data: {
     Coins: []uuid
     EncryptedOwners: map[uuid]encrypted_private_key_hex
     Metadata: map[string]string
     Receiver: Public key hex
     Redistributor: Public key hex
    }
}
Response:
  The request will fail if:
  - The redistributor is not in the list (d)
  - The coins are empty (d)
  - The encrypted data is empty (d)
  - Coins not in the encrypted owners (d)
  - The metadata is empty (d)
  - The signature is not correct (d)
  - The receiver is empty (d)
  The request will success:
  - The forward is in the database
  

Request Received:
{
   Date: time
   Signature: hex
   Action: action
   Data: {
     Hash : the sha256 hash of the coins in json
     Receiver: Public key hex
   }
}
Response:
  The request will fail if:
  - The signature is not correct
  - The hash does not exists
  - The receiver is not in the receiver of the forward
  - It is already received

Query
Request
/forwards?pub_hex=
Response
(it will return only the ones have not been received)
[]{
    Date: time
    Coins: []uuid
    EncryptedPrivateKeys: a hex of the encrypted json map[uuid]private_key_hex
}
