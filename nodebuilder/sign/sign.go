package sign

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/celestiaorg/celestia-node/libs/keystore"
	"github.com/celestiaorg/celestia-node/state"
)

var _ Module = (*API)(nil)

//go:generate mockgen -destination=mocks/api.go -package=mocks . Module
type Module interface {
	// SignPayload takes the given payload, signs it with the key by the given keyName and returns it.
	SignPayload(ctx context.Context, keyName string, payload []byte) ([]byte, error)

	// VerifySignedPayload checks whether the given payload was signed by the given publicKey.
	VerifySignedPayload(ctx context.Context, publicKey state.AccAddress, signedPayload []byte) (bool, error)
}

type signPayload struct {
	Payload   []byte `json:"payload"`
	Signature []byte `json:"signature"`
}

type module struct {
	ks keystore.Keystore
}

func newModule(ks keystore.Keystore) Module {
	return &module{
		ks: ks,
	}
}

// SignPayload takes the given payload, signs it with the key by the given keyName and returns
// the payload signed.
func (m *module) SignPayload(_ context.Context, keyName string, payload []byte) ([]byte, error) {
	signed, _, err := m.ks.Keyring().Sign(keyName, payload)
	if err != nil {
		return nil, err
	}

	sp, err := json.Marshal(signPayload{
		Payload:   payload,
		Signature: signed,
	})
	if err != nil {
		return nil, err
	}

	return sp, nil
}

// VerifySignedPayload checks whether the given payload was signed by the given publicKey.
func (m *module) VerifySignedPayload(_ context.Context, publicKey state.AccAddress, signedPayload []byte) (bool, error) {
	var sp signPayload

	err := json.Unmarshal(signedPayload, &sp)
	if err != nil {
		return false, err
	}

	err = types.VerifyAddressFormat(publicKey)
	if err != nil {
		return false, err
	}

	pkey, err := m.ks.Keyring().KeyByAddress(publicKey)
	if err != nil {
		return false, err
	}

	pk, err := pkey.GetPubKey()
	if err != nil {
		return false, err
	}

	return pk.VerifySignature(sp.Payload, sp.Signature), nil
}

// API is a wrapper around Module for the RPC.
type API struct {
	Internal struct {
		SignPayload         func(ctx context.Context, keyName string, payload []byte) ([]byte, error)                 `perm:"admin"`
		VerifySignedPayload func(ctx context.Context, publicKey state.AccAddress, signedPayload []byte) (bool, error) `perm:"admin"`
	}
}

// SignPayload takes the given payload, signs it with the key by the given keyName and returns it.
func (api *API) SignPayload(ctx context.Context, keyName string, payload []byte) ([]byte, error) {
	return api.Internal.SignPayload(ctx, keyName, payload)
}

// VerifySignedPayload checks whether the given payload was signed by the given publicKey.
func (api *API) VerifySignedPayload(ctx context.Context, publicKey state.AccAddress, signedPayload []byte) (bool, error) {
	return api.Internal.VerifySignedPayload(ctx, publicKey, signedPayload)
}
