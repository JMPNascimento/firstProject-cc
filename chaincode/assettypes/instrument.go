package assettypes

import (
	"fmt"
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
)

var Instrumento = assets.AssetType{
	Tag:         "instrumento",
	Label:       "Instrumento jurídico",
	Description: "Gerenciador de instrumentos jurídicos (contratos, petições, acordos, processos etc.)",
	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "identificador",
			Label:    "Identificador oficial",
			DataType: "numeroDoProcesso",
		},
		{
			Required: true,
			Tag:      "tipo_instrumento",
			Label:    "Tipo de instuemnto jurídico",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "status_juridico",
			Label:    "Status jurídico",
			DataType: "statusJuridico",
		},
		{
			Tag:      "advogado_responsavel",
			Label:    "Advogado responsável",
			DataType: "->advogado",
		},
		{
			Tag:      "partes",
			Label:    "Partes envolvidas",
			DataType: "[]->parte",
		},
		{
			Tag:      "objeto_assunto",
			Label:    "Objeto / Assunto",
			DataType: "string",
		},
		{
			Tag:      "data_inicio",
			Label:    "Data de início",
			DataType: "datetime",
		},
		{
			Tag:      "data_termino",
			Label:    "Data de término",
			DataType: "datetime",
		},
		{
			Tag:      "confidencialidade",
			Label:    "Status de confidencialidade (opcional)",
			DataType: "string",
		},
	},
	Validate: func(a assets.Asset) error {
		var assetMap map[string]interface{}
		switch v := interface{}(a).(type) {
		case map[string]interface{}:
			assetMap = v
		case *map[string]interface{}:
			assetMap = *v
		default:
			return fmt.Errorf("instrumento.Validate: não foi possível converter assets.Asset para map[string]interface{} (tipo: %T)", a)
		}
		diRaw, hasDi := assetMap["data_inicio"]
		dtRaw, hasDt := assetMap["data_termino"]
		if hasDi && hasDt && diRaw != nil && dtRaw != nil {
			var di time.Time
			var dt time.Time
			switch v := diRaw.(type) {
			case time.Time:
				di = v
			case string:
				parsed, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return fmt.Errorf("instrumento.Validate: data_inicio formato inválido (esperado RFC3339): %v", err)
				}
				di = parsed
			default:
				return fmt.Errorf("instrumento.Validate: data_inicio tipo inesperado %T", diRaw)
			}
			switch v := dtRaw.(type) {
			case time.Time:
				dt = v
			case string:
				parsed, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return fmt.Errorf("instrumento.Validate: data_termino formato inválido (esperado RFC3339): %v", err)
				}
				dt = parsed
			default:
				return fmt.Errorf("instrumento.Validate: data_termino tipo inesperado %T", dtRaw)
			}
			if dt.Before(di) {
				return fmt.Errorf("instrumento.Validate: data_termino não pode ser anterior a data_inicio")
			}
		}
		return nil
	},
}
