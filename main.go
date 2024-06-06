package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	bobPrivKey   = "FtSxMdxJzc4eJocNCV4kBc9QSY1H9bMAeu5qXQzocev8tMaa2uQj"
	alicePrivKey = "FsTQnqP7Kc6YdC5jfVRkUbX6wFGLh1nGUREDikCJjLRJJXGKr3ag"
	// https://github.com/btcsuite/btcd/blob/master/rpcserver.go#L77
	maxProtocolVersion = 70002
)

func buildWitnessScript(aliceKey *btcutil.WIF, bobKey *btcutil.WIF) ([32]byte, []byte, error) {

	// result hash of the game between VN and TL
	vn := sha256.Sum256([]byte("VN wins"))
	tl := sha256.Sum256([]byte("TL wins"))

	// Alice bets that VN wins
	// Bob bets that TL wins
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_SHA256)
	builder.AddOp(txscript.OP_DUP)
	builder.AddData(vn[:])
	builder.AddOp(txscript.OP_EQUAL)
	builder.AddOp(txscript.OP_IF)
	builder.AddOp(txscript.OP_DROP)
	builder.AddOp(txscript.OP_DUP)
	builder.AddOp(txscript.OP_HASH160)
	builder.AddData(btcutil.Hash160(aliceKey.SerializePubKey()))
	builder.AddOp(txscript.OP_EQUALVERIFY)
	builder.AddOp(txscript.OP_ELSE)
	builder.AddData(tl[:])
	builder.AddOp(txscript.OP_EQUALVERIFY)
	builder.AddOp(txscript.OP_DUP)
	builder.AddOp(txscript.OP_HASH160)
	builder.AddData(btcutil.Hash160(bobKey.SerializePubKey()))
	builder.AddOp(txscript.OP_EQUALVERIFY)
	builder.AddOp(txscript.OP_ENDIF)
	builder.AddOp(txscript.OP_CHECKSIG)
	pkScript, err := builder.Script()
	if err != nil {
		return [32]byte{}, []byte{}, fmt.Errorf("error build script %v", err)
	}

	witnessScriptCommitment := sha256.Sum256(pkScript)

	return witnessScriptCommitment, pkScript, nil
}

func buildFirstTx() (*wire.MsgTx, error) {
	bobKey, err := btcutil.DecodeWIF(bobPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding bob private key: %v", err)
	}

	aliceKey, err := btcutil.DecodeWIF(alicePrivKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding bob private key: %v", err)
	}

	witnessScriptCommitment, _, err := buildWitnessScript(aliceKey, bobKey)
	if err != nil {
		return nil, err
	}

	// P2WSH script
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(witnessScriptCommitment[:])
	p2wshScript, err := builder.Script()

	// random tx hash, cause we dont need to verify this tx
	txHash, err := chainhash.NewHashFromStr("aff48a9b83dc525d330ded64e1b6a9e127c99339f7246e2c89e06cd83493af9b")
	// create tx
	tx := wire.NewMsgTx(2)
	tx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  *txHash,
			Index: uint32(0),
		},
	})

	txOut := &wire.TxOut{
		Value: 490000000, PkScript: p2wshScript,
	}
	tx.AddTxOut(txOut)

	// anyone can sign here so we choose alice key
	sig, err := txscript.SignatureScript(tx, 0, []byte{}, txscript.SigHashSingle, aliceKey.PrivKey, true)
	tx.TxIn[0].SignatureScript = sig

	// log hex encoded tx
	var buf bytes.Buffer
	if err := tx.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		return nil, fmt.Errorf("failed to encode msg of type %T", tx)
	}
	fmt.Println("First tx: ", hex.EncodeToString(buf.Bytes()))

	return tx, nil
}

func buildSecondTx(preTx *wire.MsgTx) (*wire.MsgTx, error) {
	bobKey, err := btcutil.DecodeWIF(bobPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding bob private key %v", err)
	}

	aliceKey, err := btcutil.DecodeWIF(alicePrivKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding bob private key %v", err)
	}

	// build script lock
	// only alice can spent
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_DUP)
	builder.AddOp(txscript.OP_HASH160)
	builder.AddData(btcutil.Hash160(aliceKey.SerializePubKey()))
	builder.AddOp(txscript.OP_EQUALVERIFY)
	builder.AddOp(txscript.OP_CHECKSIG)
	pkScript, err := builder.Script()
	if err != nil {
		return nil, fmt.Errorf("error build script %v", err)
	}

	// create tx
	txSpent := preTx.TxOut[0]
	tx := wire.NewMsgTx(2)
	tx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  preTx.TxHash(),
			Index: uint32(0),
		},
	})

	txOut := &wire.TxOut{
		Value: 480000000, PkScript: pkScript,
	}
	tx.AddTxOut(txOut)

	inputFetcher := txscript.NewCannedPrevOutputFetcher(
		txSpent.PkScript,
		txSpent.Value,
	)
	sigHashes := txscript.NewTxSigHashes(tx, inputFetcher)

	_, witnessScript, err := buildWitnessScript(aliceKey, bobKey)
	if err != nil {
		return nil, err
	}

	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, 0, txSpent.Value, witnessScript, txscript.SigHashSingle, aliceKey.PrivKey)
	if err != nil {
		return nil, err
	}

	witness := wire.TxWitness{
		sig, aliceKey.SerializePubKey(), []byte("VN wins"), witnessScript,
	}
	tx.TxIn[0].Witness = witness

	// log hex encoded tx
	var buf bytes.Buffer
	if err := tx.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		return nil, fmt.Errorf("failed to encode msg of type %T", tx)
	}
	fmt.Println("Second tx: ", hex.EncodeToString(buf.Bytes()))

	return tx, nil
}

func main() {
	// #STEP 1: Build the first tx with a lock script that allows unlocking with the following conditions
	// - alice with message "VN wins"
	// - bob with message "TL wins"
	firstTx, err := buildFirstTx()
	if err != nil {
		fmt.Println("error build first tx ", err)
		return
	}

	// #STEP 2: Build second tx that spend first tx on alice
	secondTx, err := buildSecondTx(firstTx)
	if err != nil {
		fmt.Println("error build second tx ", err)
		return
	}

	// #STEP 3: Confirm that the above 2 tx work correctly
	blockUtxos := blockchain.NewUtxoViewpoint()
	sigCache := txscript.NewSigCache(50000)
	hashCache := txscript.NewHashCache(50000)

	inputFetcher := txscript.NewCannedPrevOutputFetcher(
		firstTx.TxOut[0].PkScript,
		firstTx.TxOut[0].Value,
	)

	blockUtxos.AddTxOut(btcutil.NewTx(firstTx), 0, 1)
	hashCache.AddSigHashes(secondTx, inputFetcher)

	// validate
	err = blockchain.ValidateTransactionScripts(
		btcutil.NewTx(secondTx), blockUtxos, txscript.StandardVerifyFlags,
		sigCache, hashCache,
	)
	if err != nil {
		fmt.Println("verify fail ", err)
	} else {
		fmt.Println("verify success")
	}

	return
}
