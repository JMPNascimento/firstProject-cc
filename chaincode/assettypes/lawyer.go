package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var Advogado = assets.AssetType{
	Tag:         "advogado",
	Label:       "Advogado",
	Description: "Advogado responsável por instrumento",
	Props: []assets.AssetProp{
		{
			Required:    true,
			IsKey:       true,
			Tag:         "oab",
			Label:       "OAB",
			Description: "Número da OAB (identificação)",
			DataType:    "string",
		},
		{
			Required: true,
			Tag:      "nome",
			Label:    "Nome",
			DataType: "string",
		},
		{
			Tag:      "contato",
			Label:    "Contato",
			DataType: "string",
		},
	},
}
