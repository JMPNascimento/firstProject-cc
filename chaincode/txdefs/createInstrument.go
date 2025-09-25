package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var CreateInstrumento = tx.Transaction{
	Tag:         "createInstrumento",
	Label:       "Create Instrumento",
	Description: "Create a new Instrumento jurídico",
	Method:      "POST",
	Callers: []accesscontrol.Caller{
		{MSP: `$org\dMSP`},
		{MSP: "orgMSP"},
	},

	Args: []tx.Argument{
		{
			Tag:      "instrumento",
			Label:    "Instrumento",
			DataType: "@object",
			Required: true,
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		instrObj, ok := req["instrumento"].(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("Argumento de instrumento inválido", 400)
		}

		instrObj["@assetType"] = "instrumento"

		if _, has := instrObj["statusJuridico"]; !has {
			instrObj["statusJuridico"] = float64(0)
		}

		instrAsset, aerr := assets.NewAsset(instrObj)
		if aerr != nil {
			return nil, errors.WrapErrorWithStatus(aerr, "Falha ao criar o instrumento", aerr.Status())
		}

		_, perr := instrAsset.PutNew(stub)
		if perr != nil {
			return nil, errors.WrapErrorWithStatus(perr, "Falha ao salvar o instrumento", perr.Status())
		}

		resp, merr := json.Marshal(instrAsset)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal response")
		}

		return resp, nil
	},
}
