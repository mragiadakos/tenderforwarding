- First we need to start ipfs daemon on its own console
$ ipfs daemon

- Secondly we need to start tendermint node on its own console, with a fresh blockchain
  The version of tendermint used is 0.21.0
$ rm -rf ~/.tendermint/
$ tendermint init
$ tendermint node

- Now we need to create an redistributor's public and private key, using the client
$ ./client g --filename=redistributor.json
The generate was successful
$ cat redistributor.json {"PrivateKey":"3081ee020100301006072a8648ce3d020106052b810400230481d63081d3020101044201241112a2a52bcfb31f68f65d3be1ec178d1f45779d9d364ad424ac30e58b1d0220e8621595831d5a07cf3b21eef44c823269dccca68d003569c4f0b3224f598de4a181890381860004014ceb7bc349fb9d49943fd9786e941f45b949b87f79f5baa288174581acbaabb4302043d447b498f74962d87cc75ea6abc1f68a960eb4686c88d160b90ef6645ea600779535284146fa179226d7f13ce810b6d10718c64b0883e484552c14186ae27facec502f4792e43b4f444bd366ac5acaed9fdd21560ffb2034e2b853c1ce9ecb60",
"PublicKey":"30819b301006072a8648ce3d020106052b810400230381860004014ceb7bc349fb9d49943fd9786e941f45b949b87f79f5baa288174581acbaabb4302043d447b498f74962d87cc75ea6abc1f68a960eb4686c88d160b90ef6645ea600779535284146fa179226d7f13ce810b6d10718c64b0883e484552c14186ae27facec502f4792e43b4f444bd366ac5acaed9fdd21560ffb2034e2b853c1ce9ecb60"}

- Create the list of redistributors in json and added in the ipfs 
$ cat redistributors.json ["30819b301006072a8648ce3d020106052b810400230381860004014ceb7bc349fb9d49943fd9786e941f45b949b87f79f5baa288174581acbaabb4302043d447b498f74962d87cc75ea6abc1f68a960eb4686c88d160b90ef6645ea600779535284146fa179226d7f13ce810b6d10718c64b0883e484552c14186ae27facec502f4792e43b4f444bd366ac5acaed9fdd21560ffb2034e2b853c1ce9ecb60"]

$ ipfs add redistributors.json 
added QmZmUEXjgZNGFbmepbRYEjG8y9SmvGc3k1SiKU92afdrht redistributors.json

- Use the IPFS hash to start the validator
$ ./server -redistributors=QmZmUEXjgZNGFbmepbRYEjG8y9SmvGc3k1SiKU92afdrht

- Create a key pair for the receiver
$ ./client g --filename=receiver.json
$ cat receiver.json 
{"PrivateKey":"3081ee020100301006072a8648ce3d020106052b810400230481d63081d3020101044201835d585921c19881f271385cbf69410912ce65d21ae4a99e0edc0b00e5f3730b2c4836412acf25981eacc40b03a9d396fa453632f2dec930310ac4de1f867d93f8a181890381860004008214502518d8d0f4fd8dc468b6083082f0c37aad9d0b69babdbd265c2c467fae2627a71563ebf5ef85072bff22ebf0c6173391e27e28e4e08acac783f15e6bf1cf01435d717a5c5242d7f4a66cdab5158e7fb9986323a204781d06b647bd119a34344d9d8cbc77604b5df99601003b5362dc5dd3737dfe5ce4e2a2ca3209b1c962309a",
"PublicKey":"30819b301006072a8648ce3d020106052b810400230381860004008214502518d8d0f4fd8dc468b6083082f0c37aad9d0b69babdbd265c2c467fae2627a71563ebf5ef85072bff22ebf0c6173391e27e28e4e08acac783f15e6bf1cf01435d717a5c5242d7f4a66cdab5158e7fb9986323a204781d06b647bd119a34344d9d8cbc77604b5df99601003b5362dc5dd3737dfe5ce4e2a2ca3209b1c962309a"}

- We have prepared the folder of the coins 
$ ll /tmp/vault
total 16
drwxr-xr-x  2 manos manos 4096 Ιουλ 15 14:16 ./
drwxrwxrwt 17 root  root  4096 Ιουλ 15 14:16 ../
-rw-r--r--  1 manos manos  226 Ιουλ 15 14:16 72b93cf8-ac6a-4d5f-9742-10da3113516c
-rw-r--r--  1 manos manos  226 Ιουλ 15 14:16 7790ffa7-6e64-417a-a542-0ceedb11e50e


- Forward the coins
$ client f --key=redistributor.json --receiver="30819b301006072a8648ce3d020106052b810400230381860004008214502518d8d0f4fd8dc468b6083082f0c37aad9d0b69babdbd265c2c467fae2627a71563ebf5ef85072bff22ebf0c6173391e27e28e4e08acac783f15e6bf1cf01435d717a5c5242d7f4a66cdab5158e7fb9986323a204781d06b647bd119a34344d9d8cbc77604b5df99601003b5362dc5dd3737dfe5ce4e2a2ca3209b1c962309a" --vault=/tmp/vault --coins=72b93cf8-ac6a-4d5f-9742-10da3113516c,7790ffa7-6e64-417a-a542-0ceedb11e50e --reason="sending money for person x to pay the y"
The forward was successfull.

- Receiver queries the forwards
$ client q --key=receiver.json
Hash:	 2712d9b5cc7a0f326e6263a7aed15ffaaa4cc904cfdcf533d670fcba3e9e3dd4
Coins:	 [72b93cf8-ac6a-4d5f-9742-10da3113516c 7790ffa7-6e64-417a-a542-0ceedb11e50e]

- Receive a hash and save the coins in a vault
$ mkdir /tmp/vaultReceiver
$ client r --key=receiver.json --vault=/tmp/vaultReceiver --hash=2712d9b5cc7a0f326e6263a7aed15ffaaa4cc904cfdcf533d670fcba3e9e3dd4
The receive was successful.


- We can can check the coins
$ ls /tmp/vaultReceiver/
72b93cf8-ac6a-4d5f-9742-10da3113516c  7790ffa7-6e64-417a-a542-0ceedb11e50e

$ cat /tmp/vaultReceiver/72b93cf8-ac6a-4d5f-9742-10da3113516c 
{"OwnerPrivateKey":"8d648034e407c4ea9d9ff6ea01caf521e1a902f26266e6205c76f3f82504540b","OwnerPublicKey":"04be911344a18f0e662136ffadf291325cd12049cd051a7fe503c7fbb63d2661","UUID":"72b93cf8-ac6a-4d5f-9742-10da3113516c","Value":1}

- And if we can try to receive again the hash, we will see that does not exists
$ client r --key=receiver.json --vault=/tmp/vaultReceiver --hash=2712d9b5cc7a0f326e6263a7aed15ffaaa4cc904cfdcf533d670fcba3e9e3dd4
Error: the hash has not been found.
