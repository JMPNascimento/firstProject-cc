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

var ExtinguishInstrumento = tx.Transaction{
	Tag:         "extinguishInstrumento",
	Label:       "Extinguir Instrumento",
	Description: "Seta status para Extinto e atualiza no histórico",
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

		instMap["status_juridico"] = float64(3)

		now := time.Now().UTC().Format(time.RFC3339)
		instMap["data_termino"] = now

		var hist []interface{}
		if h, ok := instMap["historico_status"].([]interface{}); ok {
			hist = h
		} else {
			hist = []interface{}{}
		}
		entry := map[string]interface{}{
			"status": float64(3),
			"when":   now,
			"by":     fmt.Sprint(stub.GetMSPID()),
			"motivo": req["motivo"],
		}
		hist = append(hist, entry)
		instMap["historico_status"] = hist

		updated, uerr := instAsset.Update(stub, instMap)
		if uerr != nil {
			return nil, errors.WrapErrorWithStatus(uerr, "failed to extinguish instrumento", uerr.Status())
		}

		resp, merr := json.Marshal(updated)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal response")
		}
		return resp, nil
	},
}
