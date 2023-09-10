package system

import (
	"encoding/json"
	"fmt"

	"github.com/argus-labs/starter-game-template/cardinal/tx"
	"github.com/argus-labs/starter-game-template/cardinal/utils"
	"pkg.world.dev/world-engine/cardinal/ecs"
)

func ForwardTxSystem(world *ecs.World, tq *ecs.TransactionQueue, _ *ecs.Logger) error {

	ForwardTxs := tx.ForwardTx.In(tq)

	TxTable := 

	for _, ForwardingTx := range ForwardTxs {
		world.Logger.Debug().Msg("Trying to forward")
		endpoint := ForwardingTx.Value.Endpoint
		port := ForwardingTx.Value.Port
		tx_type := ForwardingTx.Value.TxType

		target_endpoint := fmt.Sprintf("%s:%s/%s", endpoint, port, tx_type)
		tx_value := ForwardingTx.Value.TxValue
		world.Logger.Debug().Msg("tx_value")

		json_value, mar_err := json.Marshal([]byte(tx_value))
		world.Logger.Debug().Msg(string(json_value))

		world.Logger.Debug().Msg()

		if mar_err != nil {
			err_fmt := fmt.Errorf("marshaling '%s' failed as '%s': %w", tx_value, json_value, mar_err)
			world.Logger.Error().Msg(err_fmt.Error())
			world.Logger.Debug().Msg(err_fmt.Error())
			continue
		}

		_, req_err := utils.SendRequestWithJsonBody(target_endpoint, json_value, world)
		if req_err != nil {
			err_fmt := fmt.Errorf("Could not forward to '%s': %w", target_endpoint, req_err)
			world.Logger.Error().Msg(err_fmt.Error())
			world.Logger.Debug().Msg(err_fmt.Error())
			tx.ForwardTx.AddError(world, ForwardingTx.TxHash,
				err_fmt)
			continue
		}
	}

	return nil

}
