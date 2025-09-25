package assettypes

import (
	"github.com/hyperledger-labs/cc-tools/assets"
)

// Firma representa um escritório/empresa que pode referenciar vários Instrumentos
var Firma = assets.AssetType{
	Tag:         "firma",
	Label:       "Firma",
	Description: "Firma / Escritório que reúne Instrumentos",

	Props: []assets.AssetProp{
		{
			// chave primária - ajuste conforme preferir (pode ser 'cnpj' também)
			Required: true,
			IsKey:    true,
			Tag:      "nome",
			Label:    "Nome da Firma",
			DataType: "string",
			// ajuste Writers conforme sua topologia de MSPs (opcional)
			Writers: []string{`org1MSP`, "orgMSP"},
		},
		{
			Tag:      "cnpj",
			Label:    "CNPJ",
			DataType: "string",
		},
		{
			Tag:      "endereco",
			Label:    "Endereço",
			DataType: "string",
		},
		{
			Tag:      "contato",
			Label:    "Contato",
			DataType: "string",
		},
		{
			// aqui está a mudança importante: lista de referências a instrumentos
			Tag:      "instrumentos",
			Label:    "Instrumentos associados",
			DataType: "[]->instrumento",
		},
	},
}
