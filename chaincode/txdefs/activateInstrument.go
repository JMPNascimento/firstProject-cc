package txdefs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ActivateInstrumento = tx.Transaction{
	Tag:         "activateInstrumento",
	Label:       "Ativar Instrumento",
	Description: "Seta status para Em Vigência e atualiza no histórico",
	Method:      "POST",
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
			Tag:      "motivo",
			Label:    "Motivo",
			DataType: "string",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		instKey, ok := req["instrumento"].(assets.Key)
		if !ok {
			return nil, errors.NewCCError("instrumento: parâmetro inválido (deve ser ->instrumento)", 400)
		}

		instAsset, ierr := instKey.Get(stub)
		if ierr != nil {
			return nil, errors.WrapErrorWithStatus(ierr, "failed to get instrumento", ierr.Status())
		}
		instMap := (map[string]interface{})(*instAsset)

		instMap["status_juridico"] = float64(2)

		if _, exists := instMap["data_inicio"]; !exists || instMap["data_inicio"] == nil || instMap["data_inicio"] == "" {
			instMap["data_inicio"] = time.Now().UTC().Format(time.RFC3339)
		}

		var hist []interface{}
		if h, ok := instMap["historico_status"].([]interface{}); ok {
			hist = h
		} else {
			hist = []interface{}{}
		}
		now := time.Now().UTC().Format(time.RFC3339)
		entry := map[string]interface{}{
			"status": float64(2),
			"when":   now,
			"by":     fmt.Sprint(stub.GetMSPID()),
			"motivo": req["motivo"],
		}
		hist = append(hist, entry)
		instMap["historico_status"] = hist

		updated, uerr := instAsset.Update(stub, instMap)
		if uerr != nil {
			return nil, errors.WrapErrorWithStatus(uerr, "failed to activate instrumento", uerr.Status())
		}

		resp, merr := json.Marshal(updated)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal response")
		}
		return resp, nil
	},
}
