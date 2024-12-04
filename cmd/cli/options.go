package main

import "github.com/xmtp/xmtpd/pkg/config"

type GlobalOptions struct {
	Contracts config.ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Log       config.LogOptions       `group:"Log Options"       namespace:"log"`
}

type GenerateKeyOptions struct{}

type GetAllNodesOptions struct{}

type UpdateHealthOptions struct {
	AdminPrivateKey string `long:"admin-private-key" description:"Private key of the admin to administer the node"`
	NodeId          int64  `long:"node-id"           description:"NodeId to update"`
}

type UpdateAddressOptions struct {
	PrivateKey string `long:"private-key" description:"Private key of node to be updated"`
	NodeId     int64  `long:"node-id"     description:"NodeId to update"`
	Address    string `long:"address"     description:"New HTTP address"`
}

type GetPubKeyOptions struct {
	PrivateKey string `long:"private-key" description:"Private key you want the public key for" required:"true"`
}

type RegisterNodeOptions struct {
	HttpAddress     string `long:"http-address"      description:"HTTP address to register for the node"                            required:"true"`
	OwnerAddress    string `long:"owner-address"     description:"Blockchain address of the intended owner of the registration NFT" required:"true"`
	AdminPrivateKey string `long:"admin-private-key" description:"Private key of the admin to register the node"                    required:"true"`
	SigningKeyPub   string `long:"signing-key-pub"   description:"Signing key of the node to register"                              required:"true"`
}
