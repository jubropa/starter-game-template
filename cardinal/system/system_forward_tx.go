package system

import (
	"fmt"

	"github.com/argus-labs/starter-game-template/cardinal/tx"
	"github.com/argus-labs/starter-game-template/cardinal/utils"
	"pkg.world.dev/world-engine/cardinal/ecs"
)

func ForwardTxSystem(world *ecs.World, tq *ecs.TransactionQueue, _ *ecs.Logger) error {

	ForwardTxs := tx.ForwardTx.In(tq)

	for _, ForwardingTx := range ForwardTxs {
		endpoint := ForwardingTx.Value.Endpoint
		port := ForwardingTx.Value.Port
		tx_type := ForwardingTx.Value.TxType

		target_endpoint := fmt.Sprintf("%s:%s/%s", endpoint, port, tx_type)
		tx_value = ForwardingTx.Value.TxValue

		_, err := utils.SendRequestWithJsonBody(target_endpoint, tx_value)
		if err != nil {
			tx.ForwardTx.AddError(world, ForwardingTx.TxHash,
				fmt.Errorf("Could not forward to '%s': %w", target_endpoint, err))
			continue
		}

	}

	return nil

}
