package txdefs

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var SearchInstrumentosByStatus = tx.Transaction{
	Tag:         "searchInstrumentosByStatus",
	Label:       "Search Instrumentos by Status",
	Description: "Usa o query do CouchDB para encontrar instrumentos por status jurídico",
	Method:      "GET",
	Callers: []accesscontrol.Caller{
		{MSP: `$org\dMSP`},
		{MSP: "orgMSP"},
	},
	Args: []tx.Argument{
		{
			Tag:      "status",
			Label:    "Status",
			DataType: "number",
			Required: true,
		},
		{
			Tag:      "limit",
			Label:    "Limit",
			DataType: "number",
		},
	},
	ReadOnly: true,

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		statusF, ok := req["status"].(float64)
		if !ok {
			return nil, errors.NewCCError("status inválido", 400)
		}
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":      "instrumento",
				"status_juridico": statusF,
			},
		}
		// opcional limit
		if l, lok := req["limit"].(float64); lok && l > 0 {
			query["limit"] = l
		}

		results, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to execute search", 500)
		}

		out, merr := json.Marshal(results)
		if merr != nil {
			return nil, errors.WrapErrorWithStatus(merr, "failed to marshal results", 500)
		}
		return out, nil
	},
}
