// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	msg "github.com/ChainSafe/ChainBridge/message"
	"github.com/ChainSafe/log15"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

// Client is a container for all the components required to submit extrinsics
// TODO: Perhaps this would benefit an interface so we can interchange Connection and a client like this
type Client struct {
	Api     *gsrpc.SubstrateAPI
	Meta    *types.Metadata
	Genesis types.Hash
	Key     *signature.KeyringPair
}

func CreateClient(key *signature.KeyringPair, endpoint string) (*Client, error) {
	c := &Client{Key: key}
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}
	c.Api = api

	// Fetch metadata
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	c.Meta = meta

	// Fetch genesis hash
	genesisHash, err := c.Api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}
	c.Genesis = genesisHash

	return c, nil
}
func (c *Client) SetRelayerThreshold(threshold types.U32) error {
	log15.Info("Setting threshold", "threshold", threshold)
	return SubmitSudoTx(c, SetThresholdMethod, threshold)
}

func (c *Client) AddRelayer(relayer types.AccountID) error {
	log15.Info("Adding relayer", "accountId", relayer)
	return SubmitSudoTx(c, AddRelayerMethod, relayer)
}

func (c *Client) WhitelistChain(id msg.ChainId) error {
	log15.Info("Whitelisting chain", "chainId", id)
	return SubmitSudoTx(c, WhitelistChainMethod, types.U8(id))
}

func (c *Client) RegisterResource(id msg.ResourceId, method string) error {
	log15.Info("Registering resource", "rId", id, "method", []byte(method))
	return SubmitSudoTx(c, SetResourceMethod, types.NewBytes32(id), []byte(method))
}

func (c *Client) InitiateHashTransfer(hash types.Hash, destId msg.ChainId) error {
	log15.Info("Initiating hash transfer", "hash", hash.Hex())
	return SubmitTx(c, ExampleTransferHashMethod, hash, types.U8(destId))
}

func (c *Client) InitiateSubstrateNativeTransfer(amount types.U32, recipient []byte, destId msg.ChainId) error {
	log15.Info("Initiating Substrate native transfer", "amount", amount, "recipient", recipient, "destId", destId)
	return SubmitTx(c, ExampleTransferNativeMethod, amount, recipient, types.U8(destId))
}