/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package caspertest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/casper-adapter/casper"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"path/filepath"
	"testing"
)

var (
	tw *WalletManager
)

func init() {
	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "TCSPR.ini")
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	wm.WalletClient.Debug = true
	return wm
}

func TestWalletManager_GetStateRootHash(t *testing.T) {
	r, err := tw.GetStateRootHash()
	if err != nil {
		t.Errorf("GetStateRootHash failed, err: %v", err)
		return
	}
	log.Infof("stateRootHash: %s", r)
}

func TestWalletManager_GetLatestBlockInfo(t *testing.T) {
	r, err := tw.GetLatestBlockInfo()
	if err != nil {
		t.Errorf("GetLatestBlockInfo failed, err: %v", err)
		return
	}
	log.Infof("block: %v", r)
}

func TestWalletManager_GetBlockByHeight(t *testing.T) {
	block, err := tw.GetBlockByHeight(29541)
	if err != nil {
		t.Errorf("GetBlockByHeight failed, err: %v", err)
		return
	}
	log.Infof("block: %+v", block)
}

func TestWalletManager_GetBlockByHash(t *testing.T) {
	block, err := tw.GetBlockByHash("ef4a8c84dc172eb25ed3c71a9869c2a888ab6542e5b81150179e88ff55ac8621")
	if err != nil {
		t.Errorf("GetBlockByHash failed, err: %v", err)
		return
	}
	log.Infof("block: %+v", block)
}

func TestWalletManager_GetPurseUref(t *testing.T) {
	pubkey := "bf5e23418fa1b95c2be8cb900308b028d001f5626d1a2732cbe9caab71ad2cc3"
	stateRootHash, err := tw.GetStateRootHash()
	if err != nil {
		t.Errorf("GetStateRootHash failed, err: %v", err)
		return
	}
	purseUref, err := tw.GetPurseUref(stateRootHash, casper.PublicKeyToHash(pubkey, casper.ED25519_TAG))
	if err != nil {
		t.Errorf("GetPurseUref failed, err: %v", err)
		return
	}
	purseUref, _ = tw.GetPurseUref(stateRootHash, casper.PublicKeyToHash(pubkey, casper.ED25519_TAG))
	log.Infof("purseUref: %v", purseUref)
}

func TestWalletManager_GetAccountBalance(t *testing.T) {
	pubkey := "96377909058287e15ae2a3df5b77dc0abcd41136bdf8f919d5ffb412777ae475"
	stateRootHash, err := tw.GetStateRootHash()
	if err != nil {
		t.Errorf("GetStateRootHash failed, err: %v", err)
		return
	}
	r, err := tw.GetAccountBalance(stateRootHash, casper.PublicKeyToHash(pubkey, casper.ED25519_TAG))
	if err != nil {
		t.Errorf("GetAccountBalance failed, err: %v", err)
		return
	}
	log.Infof("account balance: %v", r)
}

func TestWalletManager_GetAccountBalanceByHash(t *testing.T) {
	hash := "8543cf28d54200d36842679074575ef714d5562341b8a59f0d63ad4465c11365"
	stateRootHash, err := tw.GetStateRootHash()
	if err != nil {
		t.Errorf("GetStateRootHash failed, err: %v", err)
		return
	}
	r, err := tw.GetAccountBalance(stateRootHash, hash)
	if err != nil {
		t.Errorf("GetAccountBalance failed, err: %v", err)
		return
	}
	log.Infof("account balance: %v", r)
}

func TestPublicKeyToHash(t *testing.T) {
	pubkey := "96377909058287e15ae2a3df5b77dc0abcd41136bdf8f919d5ffb412777ae475"
	hash := casper.PublicKeyToHash(pubkey, casper.ED25519_TAG)
	if len(hash) == 0 {
		t.Errorf("PublicKeyToHash failed")
		return
	}
	log.Infof("hash: %s", hash)
}

func TestWalletManager_GetBlockTransfers(t *testing.T) {
	transfers, err := tw.GetBlockTransfers(20906)
	if err != nil {
		t.Errorf("GetBlockTransfers failed, err: %v", err)
		return
	}
	for i, tx := range transfers {
		log.Infof("tx[%d]: %+v", i, tx)
	}

}

func TestWalletManager_TransferDeploy(t *testing.T) {

	privateKey, _ := hex.DecodeString("")
	senderKey := "96377909058287e15ae2a3df5b77dc0abcd41136bdf8f919d5ffb412777ae475"
	recipientKey := tw.AddressToHash("01538a824321867eebcf8506a5dc109f6cdf8abe340f2472c12d9015f24c83614c")
	//recipientKey := "8543cf28d54200d36842679074575ef714d5562341b8a59f0d63ad4465c11365"
	paymentAmount := "10000"
	transferAmount := "100000000000"
	id := "1"

	deploy, err := tw.MakeTransferDeploy(senderKey, recipientKey, transferAmount, paymentAmount, id, 0)
	if err != nil {
		t.Errorf("MakeTransferDeploy failed, err: %v", err)
		return
	}

	//serializedHeader := serializeHeader(deploy.Header)
	//log.Infof("serializedHeader Bytes: %s", hex.EncodeToString(serializedHeader))

	//deploy.hash:  499f43c30bde666764ffa503fba2630cef582cd124cbb3c96770dab07b5bf2da
	//deploy.header.bodyHash:  cc7217baa2df07b0dd1e99331a79e19290ddbb5eed8f513fec8c787a64b8e4ee

	log.Infof("deploy.hash: %s", deploy.Hash)
	log.Infof("deploy.header.bodyHash: %s", hex.EncodeToString(deploy.Header.BodyHash))

	msg, _ := hex.DecodeString(deploy.Hash)
	sig, _, ret := owcrypt.Signature(privateKey, nil, msg, owcrypt.ECC_CURVE_ED25519_NORMAL)
	if ret != owcrypt.SUCCESS {
		t.Errorf("owcrypt.Signature failed, err: %v", ret)
	}

	signature := &openwallet.KeySignature{
		EccType:   owcrypt.ECC_CURVE_ED25519,
		Address:   &openwallet.Address{
			Address:     casper.PublicKeyToHash(senderKey, casper.ED25519_TAG),
			PublicKey:   senderKey,
		},
		Signature: hex.EncodeToString(sig),
	}

	tw.AddSignatureToDeploy(deploy, signature)

	js, err := json.Marshal(deploy)
	if err != nil {
		t.Errorf("deploy to json failed, err: %v", err)
		return
	}
	log.Infof("%s", js)

	txid, err := tw.PutDeploy(deploy)
	if err != nil {
		t.Errorf("PutDeploy failed, err: %v", err)
		return
	}
	log.Infof("txid: %s", txid)
}

func TestWalletManager_GetDeployInfo(t *testing.T) {
	_, err := tw.GetDeployInfo("502c9343da4c6939173ae734993d9e8931ec9340d625fcbb6188d2e484f3cfc8")
	if err != nil {
		t.Errorf("GetDeployInfo failed, err: %v", err)
		return
	}
}

func TestWalletManager_AddressToHash(t *testing.T) {
	hash := tw.AddressToHash("01dca38c9efe06317daea1a526b36b5bd2832d71c9d69780c04f23b3a092606f91")
	log.Infof("hash: %s", hash)
}
