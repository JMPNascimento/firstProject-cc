package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var UpdateInstrumento = tx.Transaction{
	Tag:         "updateInstrumento",
	Label:       "Update Instrumento",
	Description: "Update fields of an existing Instrumento (partial update)",
	Method:      "PUT",
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
		{
			Tag:      "updates",
			Label:    "Updates",
			DataType: "@object",
			Required: true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		instrKey, ok := req["instrumento"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("O parâmetro deve ser uma referência a um asset do instrumento", 400)
		}

		instrAsset, err := instrKey.Get(stub)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Falha ao adquirir asset do instrumento", err.Status())
		}
		instrMap := (map[string]interface{})(*instrAsset)

		updates, uok := req["updates"].(map[string]interface{})
		if !uok {
			return nil, errors.NewCCError("Argumento inválido", 400)
		}

		// merge parcial: sobrescreve chaves vindas em updates
		for k, v := range updates {
			// evita sobrescrever chave/assetType
			if k == "@key" || k == "@assetType" {
				continue
			}
			instrMap[k] = v
		}

		updated, uerr := instrAsset.Update(stub, instrMap)
		if uerr != nil {
			return nil, errors.WrapError(uerr, "Falha aoatualizar asset do instrumento")
		}

		resp, merr := json.Marshal(updated)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal response")
		}
		return resp, nil
	},
}
