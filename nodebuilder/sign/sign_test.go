package sign

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/celestiaorg/celestia-node/libs/keystore"
)

const (
	keyName = "sign-test"
	// mock private key
	memorablePhrase = "mystery decrease surge raise know sense potato crisp alert crush rural music scrub flight make course maid edit often hand safe perfect dice salmon"
	// mock signed payload
	signedPayload = "7b227061796c6f6164223a226447567a6446397759586c736232466b222c227369676e6174757265223a226c43786c5857384b63514c662f74622b6f7468444159686a535946442f565a32784b616a36687266674d456967687973532f5544394834724251422b7370516c39634a2b5a6a6759434675454e6f2f4e364d417036673d3d227d"
)

func Test_module_SignPayload(t *testing.T) {
	kstore := keystore.NewMapKeystore()

	_, err := kstore.Keyring().NewAccount(keyName, memorablePhrase, sdk.GetConfig().GetFullBIP44Path(),
		keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)

	mod := newModule(kstore)

	sig, err := mod.SignPayload(context.Background(), keyName, []byte("test_payload"))
	require.NoError(t, err)

	require.Equal(t, signedPayload, hex.EncodeToString(sig))
}

func Test_module_VerifySignedPayload(t *testing.T) {
	kstore := keystore.NewMapKeystore()

	pkey, err := kstore.Keyring().NewAccount(keyName, memorablePhrase, sdk.GetConfig().GetFullBIP44Path(),
		keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)

	accAddress, err := pkey.GetAddress()

	mod := newModule(kstore)

	sig, err := hex.DecodeString(signedPayload)
	require.NoError(t, err)

	is, err := mod.VerifySignedPayload(context.Background(), accAddress, sig)
	require.NoError(t, err)
	require.True(t, is)
}
