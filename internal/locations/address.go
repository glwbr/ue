package locations

import "strings"

const (
	meterThreshold = 25.0
)

var abbreviations = map[string]string{
	"r.":   "rua",
	"av.":  "avenida",
	"dr.":  "doutor",
	"dra.": "doutora",
	"sr.":  "senhor",
	"sra.": "senhora",
	"st.":  "santo",
	"sta.": "santa",
	"vl.":  "vila",
	"jd.":  "jardim",
	"pr.":  "praca",
	"pç.":  "praca",
	"pq.":  "parque",
}

func NormalizeAddress(address string) string {
	normalized := address
	normalized = strings.TrimSpace(normalized)
	normalized = strings.ToLower(normalized)

	normalized = strings.ReplaceAll(normalized, ",", "")

	words := strings.Fields(normalized)
	for i, word := range words {
		if replacement, ok := abbreviations[word]; ok {
			words[i] = replacement
		}
	}

	normalized = strings.Join(words, " ")

	return normalized
}
