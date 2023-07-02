package util

import "strings"

func ParseStateToUF(state string) string {
	states := map[string]string{
		"acre":                "AC",
		"alagoas":             "AL",
		"amapá":               "AP",
		"amazonas":            "AM",
		"bahia":               "BA",
		"ceará":               "CE",
		"distrito federal":    "DF",
		"espírito santo":      "ES",
		"goiás":               "GO",
		"maranhão":            "MA",
		"mato grosso":         "MT",
		"mato grosso do sul":  "MS",
		"minas gerais":        "MG",
		"pará":                "PA",
		"paraíba":             "PB",
		"paraná":              "PR",
		"pernambuco":          "PE",
		"piauí":               "PI",
		"rio de janeiro":      "RJ",
		"rio grande do norte": "RN",
		"rio grande do sul":   "RS",
		"rondônia":            "RO",
		"roraima":             "RR",
		"santa catarina":      "SC",
		"são paulo":           "SP",
		"sergipe":             "SE",
		"tocantins":           "TO",
	}

	uf, ok := states[strings.ToLower(state)]
	if !ok {
		return state
	}

	return uf
}
