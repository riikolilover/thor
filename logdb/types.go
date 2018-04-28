package logdb

import (
	"math/big"

	"github.com/vechain/thor/block"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
)

//Event represents tx.Event that can be stored in db.
type Event struct {
	BlockID     thor.Bytes32
	Index       uint32
	BlockNumber uint32
	BlockTime   uint64
	TxID        thor.Bytes32
	TxOrigin    thor.Address //contract caller
	Address     thor.Address // always a contract address
	Topics      [5]*thor.Bytes32
	Data        []byte
}

//newEvent converts tx.Event to Event.
func newEvent(header *block.Header, index uint32, txID thor.Bytes32, txOrigin thor.Address, txEvent *tx.Event) *Event {
	ev := &Event{
		BlockID:     header.ID(),
		Index:       index,
		BlockNumber: header.Number(),
		BlockTime:   header.Timestamp(),
		TxID:        txID,
		TxOrigin:    txOrigin,
		Address:     txEvent.Address, // always a contract address
		Data:        txEvent.Data,
	}
	for i := 0; i < len(txEvent.Topics) && i < len(ev.Topics); i++ {
		ev.Topics[i] = &txEvent.Topics[i]
	}
	return ev
}

//Transfer represents tx.Transfer that can be stored in db.
type Transfer struct {
	BlockID     thor.Bytes32
	Index       uint32
	BlockNumber uint32
	BlockTime   uint64
	TxID        thor.Bytes32
	TxOrigin    thor.Address
	From        thor.Address
	To          thor.Address
	Value       *big.Int
}

//newTransfer converts tx.Transfer to Transfer.
func newTransfer(header *block.Header, index uint32, txID thor.Bytes32, txOrigin thor.Address, transfer *tx.Transfer) *Transfer {
	return &Transfer{
		BlockID:     header.ID(),
		Index:       index,
		BlockNumber: header.Number(),
		BlockTime:   header.Timestamp(),
		TxID:        txID,
		TxOrigin:    txOrigin,
		From:        transfer.Sender,
		To:          transfer.Recipient,
		Value:       transfer.Amount,
	}
}

type RangeType string

const (
	Block RangeType = "block"
	Time  RangeType = "time"
)

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type Range struct {
	Unit RangeType
	From uint64
	To   uint64
}

type Options struct {
	Offset uint64
	Limit  uint64
}

//EventFilter filter
type EventFilter struct {
	Address  *thor.Address // always a contract address
	TopicSet [][5]*thor.Bytes32
	Range    *Range
	Options  *Options
	Order    Order //default asc
}

type AddressSet struct {
	TxOrigin *thor.Address //who send transaction
	From     *thor.Address //who transferred tokens
	To       *thor.Address //who recieved tokens
}

type TransferFilter struct {
	TxID        *thor.Bytes32
	AddressSets []*AddressSet
	Range       *Range
	Options     *Options
	Order       Order //default asc
}
