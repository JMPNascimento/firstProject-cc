package datatypes

import (
	"fmt"
	"strconv"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

type StatusJuridicoType float64

const (
	StatusEmTramitacao StatusJuridicoType = iota
	StatusArquivado
	StatusEmVigencia
	StatusExtinto
)

func (s StatusJuridicoType) CheckType() errors.ICCError {
	switch s {
	case StatusEmTramitacao, StatusArquivado, StatusEmVigencia, StatusExtinto:
		return nil
	default:
		return errors.NewCCError("status_juridico: valor inválido", 400)
	}
}

var StatusJuridico = assets.DataType{
	AcceptedFormats: []string{"number"},
	DropDownValues: map[string]interface{}{
		"Em Tramitação": StatusEmTramitacao,
		"Arquivado":     StatusArquivado,
		"Em Vigência":   StatusEmVigencia,
		"Extinto":       StatusExtinto,
	},
	Description: `Status jurídico do instrumento (enum)`,

	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		var dataVal float64

		switch v := data.(type) {
		case float64:
			dataVal = v
		case int:
			dataVal = float64(v)
		case StatusJuridicoType:
			dataVal = float64(v)
		case string:
			var err error
			dataVal, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return "", nil, errors.WrapErrorWithStatus(err, "status_juridico: valor deve ser inteiro/number", 400)
			}
		default:
			return "", nil, errors.NewCCError("status_juridico: asset property must be an integer/number", 400)
		}

		retVal := StatusJuridicoType(dataVal)
		err := retVal.CheckType()
		return fmt.Sprint(retVal), retVal, err
	},
}
