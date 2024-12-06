package sign

import (
	"go.uber.org/fx"

	"github.com/celestiaorg/celestia-node/libs/keystore"
	"github.com/celestiaorg/celestia-node/nodebuilder/node"
)

// ConstructModule constructs a new module for the given node type.
//
// The module (sign) is constructed using the provided keystore for node type
// node.Light, node.Full, node.Bridge.
func ConstructModule(tp node.Type) fx.Option {
	switch tp {
	case node.Light, node.Full, node.Bridge:
		return fx.Module(
			"sign",
			fx.Invoke(fx.Provide(func(ks keystore.Keystore) Module {
				return newModule(ks)
			})),
		)
	default:
		panic("invalid node type")
	}
}
