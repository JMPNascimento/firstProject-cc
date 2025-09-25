package txdefs

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

// GetInstrumentHistory - retorna histórico de versões do asset (GetHistoryForKey)
// GET Method
var GetInstrumentHistory = tx.Transaction{
	Tag:         "getInstrumentHistory",
	Label:       "Get Instrument History",
	Description: "Return the history of an Instrumento using history database",
	Method:      "GET",
	Callers: []accesscontrol.Caller{
		{MSP: `$org\dMSP`},
		{MSP: "orgMSP"},
	},
	Args: []tx.Argument{
		{
			Tag:      "instrumento",
			Label:    "Instrumento",
			DataType: "->instrumento",
			Required: true,
		},
	},
	ReadOnly: true,

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		instKey, ok := req["instrumento"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("instrumento: parâmetro inválido", 400)
		}

		// obter map atual para recuperar a chave string (campo @key)
		instMap, err := instKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to get instrumento", err.Status())
		}
		keyStr, ok := instMap["@key"].(string)
		if !ok {
			return nil, errors.NewCCError("instrumento: chave interna (@key) não encontrada", 500)
		}

		historicIter, herr := stub.Stub.GetHistoryForKey(keyStr)
		if herr != nil {
			return nil, errors.WrapError(herr, "failed to get history from ledger")
		}
		defer historicIter.Close()

		var history []map[string]interface{}
		for historicIter.HasNext() {
			mod, nerr := historicIter.Next()
			if nerr != nil {
				return nil, errors.WrapError(nerr, "error iterating history")
			}

			var value map[string]interface{}
			_ = json.Unmarshal(mod.Value, &value)

			// Convert timestamp to RFC3339 if possible
			var ts interface{}
			if mod.Timestamp != nil {
				t := time.Unix(mod.Timestamp.Seconds, int64(mod.Timestamp.Nanos)).UTC().Format(time.RFC3339)
				ts = t
			} else {
				ts = nil
			}

			history = append(history, map[string]interface{}{
				"txid":      mod.TxId,
				"timestamp": ts,
				"value":     value,
				"isDelete":  mod.IsDelete,
			})
		}

		out, merr := json.Marshal(history)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal history")
		}
		return out, nil
	},
}
