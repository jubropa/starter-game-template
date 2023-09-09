package system

import (
	"fmt"

	comp "github.com/argus-labs/starter-game-template/cardinal/component"
	"github.com/argus-labs/starter-game-template/cardinal/tx"
	"github.com/argus-labs/starter-game-template/cardinal/utils"
	"pkg.world.dev/world-engine/cardinal/ecs"
	"pkg.world.dev/world-engine/cardinal/ecs/filter"
	"pkg.world.dev/world-engine/cardinal/ecs/storage"
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

func AttackSystem(world *ecs.World, tq *ecs.TransactionQueue, _ *ecs.Logger) error {

	attackTx := tx.AttackPlayerMsg{
		TargetNickname: "cool guy",
	}

	// 	http call using url and attackTx

	// Get all the transactions that are of type CreatePlayer from the tx queue
	attackTxs := tx.AttackPlayer.In(tq)

	// Create an index of player tags to its health component
	playerTagToID := map[string]storage.EntityID{}
	ecs.NewQuery(filter.Exact(comp.Player, comp.Health)).Each(world, func(id storage.EntityID) bool {
		player, err := comp.Player.Get(world, id)
		if err != nil {
			return true
		}

		playerTagToID[player.Nickname] = id
		return true
	})

	// Iterate through all transactions and process them individually.
	// DEV: it's important here that you don't break out of the loop or return an error here
	// or otherwise the rest of the transaction will not be processed & get dropped.
	// In the future, you will be able to add error receipts to transaction receipts.
	for _, attack := range attackTxs {
		target := attack.Value.TargetNickname
		targetPlayerID, ok := playerTagToID[target]
		// If the target player doesn't exist, skip this transaction
		if !ok {
			tx.AttackPlayer.AddError(world, attack.TxHash,
				fmt.Errorf("target %q does not exist", target))
			continue
		}

		// Get the health component for the target player
		health, err := comp.Health.Get(world, targetPlayerID)
		if err != nil {
			tx.AttackPlayer.AddError(world, attack.TxHash,
				fmt.Errorf("can't get health for %q: %w", target, err))
			continue
		}

		// Inflict damage and update the component
		health.HP -= 10
		if err := comp.Health.Set(world, targetPlayerID, health); err != nil {
			tx.AttackPlayer.AddError(world, attack.TxHash,
				fmt.Errorf("failed to set health on %q: %w", target, err))
			continue
		}
	}

	return nil
}
