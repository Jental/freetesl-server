package enums

type CardKeyword byte

const (
	CardKeywordBreakthrough CardKeyword = 1
	CardKeywordCharge       CardKeyword = 2
	CardKeywordDrain        CardKeyword = 3
	CardKeywordGuard        CardKeyword = 4
	CardKeywordLethal       CardKeyword = 5
	CardKeywordMobilize     CardKeyword = 6
	CardKeywordRally        CardKeyword = 7
	CardKeywordRegenerate   CardKeyword = 8
	CardKeywordWard         CardKeyword = 9
)

var CardKeywordName = map[CardKeyword]string{
	CardKeywordBreakthrough: "Breakthrough",
	CardKeywordCharge:       "Charge",
	CardKeywordDrain:        "Drain",
	CardKeywordGuard:        "Guard",
	CardKeywordLethal:       "Lethal",
	CardKeywordMobilize:     "Mobilize",
	CardKeywordRally:        "Rally",
	CardKeywordRegenerate:   "Regenerate",
	CardKeywordWard:         "Ward",
}
