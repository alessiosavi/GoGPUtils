// Package stopwords provides language-specific stopword sets for use with
// textnorm.RemoveStopwords. The sets are exported as map[string]struct{} so
// callers get O(1) membership checks and can pass them straight to the
// RemoveStopwords stage.
//
// The variables are documented as read-only — Go has no const map, but
// callers MUST NOT mutate them. Use Union to combine sets safely.
package stopwords

// English is a default English stopword set covering ~150 common function
// words. Suitable for search-key filtering, embeddings hygiene, and slug
// generation.
var English = map[string]struct{}{
	"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {},
	"all": {}, "am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "as": {},
	"at": {}, "be": {}, "because": {}, "been": {}, "before": {}, "being": {},
	"below": {}, "between": {}, "both": {}, "but": {}, "by": {}, "can": {},
	"could": {}, "did": {}, "do": {}, "does": {}, "doing": {}, "don": {},
	"down": {}, "during": {}, "each": {}, "few": {}, "for": {}, "from": {},
	"further": {}, "had": {}, "has": {}, "have": {}, "having": {}, "he": {},
	"her": {}, "here": {}, "hers": {}, "herself": {}, "him": {}, "himself": {},
	"his": {}, "how": {}, "i": {}, "if": {}, "in": {}, "into": {}, "is": {},
	"it": {}, "its": {}, "itself": {}, "just": {}, "me": {}, "more": {},
	"most": {}, "my": {}, "myself": {}, "no": {}, "nor": {}, "not": {},
	"now": {}, "of": {}, "off": {}, "on": {}, "once": {}, "only": {}, "or": {},
	"other": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {}, "over": {},
	"own": {}, "same": {}, "she": {}, "should": {}, "so": {}, "some": {},
	"such": {}, "than": {}, "that": {}, "the": {}, "their": {}, "theirs": {},
	"them": {}, "themselves": {}, "then": {}, "there": {}, "these": {},
	"they": {}, "this": {}, "those": {}, "through": {}, "to": {}, "too": {},
	"under": {}, "until": {}, "up": {}, "very": {}, "was": {}, "we": {},
	"were": {}, "what": {}, "when": {}, "where": {}, "which": {}, "while": {},
	"who": {}, "whom": {}, "why": {}, "will": {}, "with": {}, "would": {},
	"you": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
}

// French is a starter French stopword set. Expected to grow; treat the API
// (variable name + type) as the contract, the contents as data.
var French = map[string]struct{}{
	"à": {}, "au": {}, "aux": {}, "avec": {}, "ce": {}, "ces": {}, "cette": {},
	"d": {}, "dans": {}, "de": {}, "des": {}, "du": {}, "elle": {}, "en": {},
	"est": {}, "et": {}, "eu": {}, "il": {}, "ils": {}, "j": {}, "je": {},
	"l": {}, "la": {}, "le": {}, "les": {}, "leur": {}, "lui": {}, "ma": {},
	"mais": {}, "me": {}, "mes": {}, "moi": {}, "mon": {}, "n": {}, "ne": {},
	"nos": {}, "notre": {}, "nous": {}, "on": {}, "ou": {}, "par": {},
	"pas": {}, "pour": {}, "qu": {}, "que": {}, "qui": {}, "s": {}, "sa": {},
	"se": {}, "ses": {}, "son": {}, "sur": {}, "ta": {}, "te": {}, "tes": {},
	"toi": {}, "ton": {}, "tu": {}, "un": {}, "une": {}, "vos": {},
	"votre": {}, "vous": {}, "y": {},
}

// Italian is a starter Italian stopword set. Expected to grow; treat the API
// (variable name + type) as the contract, the contents as data.
var Italian = map[string]struct{}{
	"a": {}, "ad": {}, "al": {}, "alla": {}, "alle": {}, "agli": {}, "ai": {},
	"allo": {}, "anche": {}, "che": {}, "chi": {}, "ci": {}, "come": {},
	"con": {}, "da": {}, "dal": {}, "dalla": {}, "del": {}, "della": {},
	"delle": {}, "dello": {}, "degli": {}, "dei": {}, "di": {}, "e": {},
	"ed": {}, "è": {}, "gli": {}, "i": {}, "il": {}, "in": {}, "io": {},
	"l": {}, "la": {}, "le": {}, "lei": {}, "li": {}, "lo": {}, "loro": {},
	"lui": {}, "ma": {}, "me": {}, "mi": {}, "mio": {}, "ne": {}, "nei": {},
	"nel": {}, "nella": {}, "nelle": {}, "no": {}, "noi": {}, "non": {},
	"o": {}, "per": {}, "più": {}, "questa": {}, "queste": {}, "questi": {},
	"questo": {}, "se": {}, "sei": {}, "si": {}, "sia": {}, "sono": {},
	"su": {}, "sul": {}, "sulla": {}, "te": {}, "ti": {}, "tra": {}, "tu": {},
	"tuo": {}, "un": {}, "una": {}, "uno": {}, "voi": {},
}

// Union returns a new set containing every key found in any of the input
// sets. It does not mutate any input. With zero inputs it returns a non-nil
// empty map.
func Union(sets ...map[string]struct{}) map[string]struct{} {
	total := 0
	for _, s := range sets {
		total += len(s)
	}
	out := make(map[string]struct{}, total)
	for _, s := range sets {
		for k := range s {
			out[k] = struct{}{}
		}
	}
	return out
}
