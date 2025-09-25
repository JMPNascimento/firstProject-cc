package txdefs

import (
	"encoding/json"
	"time"

	"fmt"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
)

var ArchiveInstrumento = tx.Transaction{
	Tag:         "archiveInstrumento",
	Label:       "Archive Instrumento",
	Description: "Seta status para Arquivado e atualiza no histórico",
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
			return nil, errors.NewCCError("instrumento: parâmetro inválido", 400)
		}

		instAsset, err := instKey.Get(stub)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to get instrumento", err.Status())
		}
		instMap := (map[string]interface{})(*instAsset)

		// atualiza status para "Arquivado" (assumimos enum 1)
		instMap["status_juridico"] = float64(1)

		// adiciona entrada no histórico interno (campo historico_status)
		var hist []interface{}
		if h, ok := instMap["historico_status"].([]interface{}); ok {
			hist = h
		}
		entry := map[string]interface{}{
			"status": float64(1),
			"when":   time.Now().UTC().Format(time.RFC3339),
			"by":     fmt.Sprint(stub.GetMSPID()), // usa método do stubwrapper (comum no cc-tools)
			"motivo": req["motivo"],
		}
		hist = append(hist, entry)
		instMap["historico_status"] = hist

		// grava (Update faz validações)
		updated, uerr := instAsset.Update(stub, instMap)
		if uerr != nil {
			return nil, errors.WrapError(uerr, "failed to archive instrumento")
		}

		resp, merr := json.Marshal(updated)
		if merr != nil {
			return nil, errors.WrapError(merr, "failed to marshal response")
		}
		return resp, nil
	},
}
