package genesis

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vechain/thor/block"
	cs "github.com/vechain/thor/contracts"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/thor"
)

var Dev = &dev{}

type dev struct {
}

func (d *dev) Accounts() (accs [10]struct {
	Address    thor.Address
	PrivateKey *ecdsa.PrivateKey
}) {

	privKeys := [10]string{
		"dce1443bd2ef0c2631adc1c67e5c93f13dc23a41c18b536effbbdcbcdb96fb65",
		"321d6443bc6177273b5abf54210fe806d451d6b7973bccc2384ef78bbcd0bf51",
		"2d7c882bad2a01105e36dda3646693bc1aaaa45b0ed63fb0ce23c060294f3af2",
		"593537225b037191d322c3b1df585fb1e5100811b71a6f7fc7e29cca1333483e",
		"ca7b25fc980c759df5f3ce17a3d881d6e19a38e651fc4315fc08917edab41058",
		"88d2d80b12b92feaa0da6d62309463d20408157723f2d7e799b6a74ead9a673b",
		"fbb9e7ba5fe9969a71c6599052237b91adeb1e5fc0c96727b66e56ff5d02f9d0",
		"547fb081e73dc2e22b4aae5c60e2970b008ac4fc3073aebc27d41ace9c4f53e9",
		"c8c53657e41a8d669349fc287f57457bd746cb1fcfc38cf94d235deb2cfca81b",
		"87e0eba9c86c494d98353800571089f316740b0cb84c9a7cdf2fe5c9997c7966",
	}
	for i, str := range privKeys {
		pk, err := crypto.HexToECDSA(str)
		if err != nil {
			panic(err)
		}
		addr := crypto.PubkeyToAddress(pk.PublicKey)
		accs[i].Address = thor.Address(addr)
		accs[i].PrivateKey = pk
	}
	return
}

func (d *dev) Build(stateCreator *state.Creator) (*block.Block, error) {
	builder := new(Builder).
		GasLimit(thor.InitialGasLimit).
		Timestamp(1516333644).
		Alloc(cs.Authority.Address, &big.Int{}, cs.Authority.RuntimeBytecodes()).
		Alloc(cs.Energy.Address, &big.Int{}, cs.Energy.RuntimeBytecodes()).
		Alloc(cs.Params.Address, &big.Int{}, cs.Params.RuntimeBytecodes()).
		Call(cs.Authority.PackInitialize()).
		Call(cs.Energy.PackInitialize(cs.Voting.Address)).
		Call(cs.Params.PackInitialize()).
		Call(cs.Params.PackSet(cs.ParamRewardRatio, big.NewInt(3e17))).
		Call(cs.Params.PackSet(cs.ParamBaseGasPrice, big.NewInt(1000)))

	for _, a := range d.Accounts() {
		balance, _ := new(big.Int).SetString("10000000000000000000000", 10)
		builder.Alloc(a.Address, balance, nil).
			Call(cs.Authority.PackAuthorize(a.Address, "a1")).
			Call(cs.Energy.PackCharge(a.Address, balance))
	}

	genesis, err := builder.Build(stateCreator)
	if err != nil {
		return nil, err
	}
	return genesis.WithSignature([]byte{2}), nil
}
