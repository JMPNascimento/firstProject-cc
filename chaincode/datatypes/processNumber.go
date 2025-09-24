package datatypes

import (
	"regexp"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

var numeroDoProcessoREGEX = regexp.MustCompile(`^[A-Za-z0-9\-/\.]{3,60}$`)

var numeroDoProcesso = assets.DataType{
	AcceptedFormats: []string{"string"},
	Description:     "Identificador oficial de um instrumento",
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		var s string

		switch v := data.(type) {
		case string:
			s = v
		case []byte:
			s = string(v)
		default:
			return "", nil, errors.NewCCError("O identificador do processo deve ser uma string", 400)
		}
		if !numeroDoProcessoREGEX.MatchString(s) {
			return "", nil, errors.NewCCError("Formato de identifacdor de processo inválido", 400)
		}
		return s, s, nil
	},
}

var _ = numeroDoProcessoREGEX
var _ = numeroDoProcesso
