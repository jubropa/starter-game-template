package tx

import (
	"pkg.world.dev/world-engine/cardinal/ecs"
)

type ForwardTxMsg struct {
	Endpoint       string `json:"endpoint"`
	Port           string `json:"port"`
	TxTypeStruct   any    `json:"tx_type"`
	TxTypeEndpoint string `json:"tx_type_endpoint"`
	TxValue        string `json:"tx_value"`
}

type ForwardTxMsgReply struct {
	Endpoint string `json:"endpoint"`
	Port     string `json:"port"`
	TxType   string `json:"tx_type"`
	TxResp   string `json:"tx_resp"`
}

var ForwardTx = ecs.NewTransactionType[ForwardTxMsg, ForwardTxMsgReply]("forward-tx")
