package models

type ForwardModel struct {
	Coins           []string
	EncryptedOwners map[string]string
	Metadata        map[string]string
	Receiver        string
	Redistributor   string
}
