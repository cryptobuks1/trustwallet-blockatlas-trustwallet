package tezos

import (
	"fmt"
	"time"

	"github.com/trustwallet/golibs/types"
)

const (
	TxTypeTransaction string = "transaction"
	TxTypeDelegation  string = "delegation"

	TxStatusApplied string = "applied"
)

type (
	Account struct {
		Balance  string `json:"balance"`
		Delegate string `json:"delegate"`
	}

	AccountBalance struct {
		Balance string `json:"balance"`
	}

	ExplorerAccount struct {
		Transactions []Transaction `json:"ops"`
	}

	Transaction struct {
		Delegate  string  `json:"delegate"` // Current delegate (may be self when registered as delegate).
		Errors    []Error `json:"errors"`   // Operation status applied, failed, backtracked, skipped.
		Fee       float64 `json:"fee"`      // Total fee paid (and frozen) by all operations.
		Hash      string  `json:"hash"`     // Operation hash.
		Height    uint64  `json:"height"`
		IsSuccess bool    `json:"is_success"` // Flag indicating operation was successfully applied.
		Receiver  string  `json:"receiver"`
		Sender    string  `json:"sender"`
		Stat      string  `json:"status"` // Operation status applied, failed, backtracked, skipped.
		Time      string  `json:"time"`   // Block time at which the operation was included on-chain e.g: 2019-09-28T13:10:51Z
		Type      string  `json:"type"`   // Operation type, one of activate_account, double_baking_evidence, double_endorsement_evidence, seed_nonce_revelation, transaction, origination, delegation, reveal, endorsement, proposals, ballot.
		Volume    float64 `json:"volume"`
	}

	Error struct {
		ID   string `json:"id"`
		Kind string `json:"kind"`
	}

	Validator struct {
		Address string `json:"pkh"`
	}

	ActivityValidatorInfo struct {
		Deactivated bool `json:"deactivated"`
	}

	Baker struct {
		Address           string  `json:"address"`
		Name              string  `json:"name"`
		Logo              string  `json:"logo"`
		FreeSpace         float64 `json:"freeSpace"`
		Fee               float64 `json:"fee"`
		MinDelegation     float64 `json:"minDelegation"`
		OpenForDelegation bool    `json:"openForDelegation"`
		EstimatedRoi      float64 `json:"estimatedRoi"`
		ServiceHealth     string  `json:"serviceHealth"`
	}
)

func (t *Transaction) Status() types.Status {
	switch t.Stat {
	case TxStatusApplied:
		return types.StatusCompleted
	default:
		return types.StatusError
	}
}

func (t *Transaction) ErrorMsg() string {
	if !t.IsSuccess && len(t.Errors) > 0 {
		return fmt.Sprintf("%s %s", t.Errors[0].ID, t.Errors[0].Kind)
	} else {
		return ""
	}
}

func (t *Transaction) Title(address string) (types.KeyTitle, bool) {
	if t.Type == TxTypeDelegation {
		if address == t.Sender && t.Delegate != "" && t.Receiver == "" {
			return types.AnyActionDelegation, true
		}

		if address == t.Sender && t.Delegate == "" && t.Receiver != "" {
			return types.AnyActionUndelegation, true
		}
	}

	return "unsupported title", false
}

func (t *Transaction) BlockTimestamp() int64 {
	unix := int64(0)
	date, err := time.Parse(time.RFC3339, t.Time)
	if err == nil {
		unix = date.Unix()
	}
	return unix
}

func (t *Transaction) TransferType() (types.TransactionType, bool) {
	switch t.Type {
	case TxTypeTransaction:
		return types.TxTransfer, true
	case TxTypeDelegation:
		return types.TxAnyAction, true
	default:
		return "unsupported type", false
	}
}

func (t *Transaction) Direction(address string) types.Direction {
	if t.Sender == address && t.Receiver == address {
		return types.DirectionSelf
	}
	if t.Sender == address && t.Receiver != address {
		return types.DirectionOutgoing
	}

	return types.DirectionIncoming
}

func (t *Transaction) GetReceiver() string {
	if t.Receiver != "" {
		return t.Receiver
	} else {
		return t.Delegate
	}
}
