package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

var Parte = assets.AssetType{
	Tag:         "parte",
	Label:       "Parte envolvida",
	Description: "Pessoa ou entidade envolvida no instrumento",
	Props: []assets.AssetProp{
		{
			Required:    true,
			IsKey:       true,
			Tag:         "id",
			Label:       "ID",
			Description: "Identificador (CPF/CNPJ)",
			DataType:    "string",
		},
		{
			Required: true,
			Tag:      "nome",
			Label:    "Nome / Razão Social",
			DataType: "string",
		},
		{
			Tag:      "contato",
			Label:    "Contato",
			DataType: "string",
		},
	},
}
