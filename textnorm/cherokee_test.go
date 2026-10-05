package textnorm

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/cases"
)

// All 86 Cherokee pairs, in matching order, from Unicode 17 CaseFolding.txt:
// https://www.unicode.org/Public/17.0.0/ucd/CaseFolding.txt
// Listed lowercase characters map to uppercase; unlisted uppercase map to self.
const cherokeeLower = "ꭰꭱꭲꭳꭴꭵꭶꭷꭸꭹꭺꭻꭼꭽꭾꭿꮀꮁꮂꮃꮄꮅꮆꮇꮈꮉꮊꮋꮌꮍꮎꮏꮐꮑꮒꮓꮔꮕꮖꮗꮘꮙꮚꮛꮜꮝꮞꮟꮠꮡꮢꮣꮤꮥꮦꮧꮨꮩꮪꮫꮬꮭꮮꮯꮰꮱꮲꮳꮴꮵꮶꮷꮸꮹꮺꮻꮼꮽꮾꮿᏸᏹᏺᏻᏼᏽ"
const cherokeeUpper = "ᎠᎡᎢᎣᎤᎥᎦᎧᎨᎩᎪᎫᎬᎭᎮᎯᎰᎱᎲᎳᎴᎵᎶᎷᎸᎹᎺᎻᎼᎽᎾᎿᏀᏁᏂᏃᏄᏅᏆᏇᏈᏉᏊᏋᏌᏍᏎᏏᏐᏑᏒᏓᏔᏕᏖᏗᏘᏙᏚᏛᏜᏝᏞᏟᏠᏡᏢᏣᏤᏥᏦᏧᏨᏩᏪᏫᏬᏭᏮᏯᏰᏱᏲᏳᏴᏵ"

// Keep the oracle independent of the production range mapping, and preserve
// malformed bytes by replacing only the explicitly listed UTF-8 sequences.
var cherokeeFoldOracle = func() *strings.Replacer {
	lower, upper := []rune(cherokeeLower), []rune(cherokeeUpper)
	pairs := make([]string, 0, 2*len(lower))
	for i, r := range lower {
		pairs = append(pairs, string(r), string(upper[i]))
	}
	return strings.NewReplacer(pairs...)
}()

func cherokeePipelines() []struct {
	name string
	pipe Pipeline
} {
	return []struct {
		name string
		pipe Pipeline
	}{
		{"FoldCase", New().FoldCase()},
		{"CanonicalPreset", CanonicalPreset()},
		{"MeaningPreset", MeaningPreset()},
		{"SearchPreset", SearchPreset()},
		{"CanonicalPreset/width", CanonicalPreset(WithWidthFold())},
		{"MeaningPreset/width", MeaningPreset(WithWidthFold())},
		{"SearchPreset/width", SearchPreset(WithWidthFold())},
	}
}

func checkCherokeeFold(t *testing.T, p Pipeline, in, want string) {
	t.Helper()
	got, err := p.Run(in)
	if err != nil || got != want {
		t.Fatalf("Run(%q) = %q, %v; want %q, nil", in, got, err, want)
	}
	again, err := p.Run(got)
	if err != nil || again != got {
		t.Fatalf("second Run(%q) = %q, %v; want %q, nil", got, again, err, got)
	}
}

func TestCherokeeCasePairs(t *testing.T) {
	lower, upper := []rune(cherokeeLower), []rune(cherokeeUpper)
	if len(lower) != 86 || len(upper) != 86 {
		t.Fatal("fixture must cover all 86 Cherokee case pairs")
	}
	for _, p := range cherokeePipelines() {
		t.Run(p.name, func(t *testing.T) {
			for i, lo := range lower {
				t.Run(fmt.Sprintf("%U", lo), func(t *testing.T) {
					want := string(upper[i])
					checkCherokeeFold(t, p.pipe, string(lo), want)
					checkCherokeeFold(t, p.pipe, want, want)
				})
			}
		})
	}
}

func TestCherokeeFoldBoundaries(t *testing.T) {
	for _, p := range cherokeePipelines() {
		for _, n := range []int{0, 10, 24, 27, 127, 128, 129, 255, 256, 257, 4095, 4096, 4097} {
			t.Run(fmt.Sprintf("%s/prefix=%d", p.name, n), func(t *testing.T) {
				prefix := strings.Repeat("a", n)
				in := prefix + cherokeeLower + "z" + cherokeeUpper + "0"
				want := prefix + cherokeeUpper + "z" + cherokeeUpper + "0"
				checkCherokeeFold(t, p.pipe, in, want)
			})
		}
	}
}

func TestCherokeeFoldMixedText(t *testing.T) {
	for _, tc := range []struct {
		name     string
		pipe     Pipeline
		in, want string
	}{
		{"FoldCase", New().FoldCase(), "\tStraße ꮳ ᏸ CAFÉ Σς ﬁ 中文 😀\n", "\tstrasse Ꮳ Ᏸ café σσ fi 中文 😀\n"},
		{"CanonicalPreset", CanonicalPreset(), "\tStraße ꮳ ᏸ CAFÉ Σς ﬁ 中文 😀\n", "strasse Ꮳ Ᏸ cafe σσ fi 中文 😀"},
		{"MeaningPreset", MeaningPreset(), "\tStraße ꮳ ᏸ CAFÉ Σς ﬁ 中文 😀\n", "strasse Ꮳ Ᏸ cafe σσ fi 中文"},
		{"SearchPreset", SearchPreset(), "\tStraße ꮳ ᏸ CAFÉ Σς ﬁ 中文 😀\n", "strasse Ꮳ Ᏸ cafe σσ fi 中文"},
		{"FoldCase/invalid", New().FoldCase(), "\xffStraßeᏣ\x80ᏸ\x00中文\xed\xa0\x80", "\xffstrasseᏣ\x80Ᏸ\x00中文\xed\xa0\x80"},
		{"CanonicalPreset/invalid", CanonicalPreset(), "\xffStraßeᏣ\x80ᏸ\x00中文\xed\xa0\x80", "�strasseᏣ�Ᏸ中文���"},
		{"MeaningPreset/invalid", MeaningPreset(), "\xffStraßeᏣ\x80ᏸ\x00中文\xed\xa0\x80", "strasseᏣ Ᏸ中文"},
		{"SearchPreset/invalid", SearchPreset(), "\xffStraßeᏣ\x80ᏸ\x00中文\xed\xa0\x80", "strasseᏣᏰ中文"},
		{"saved/Cherokee", CanonicalPreset(), "00000000000000000000000000000ꮒ0", "00000000000000000000000000000Ꮒ0"},
	} {
		t.Run(tc.name, func(t *testing.T) { checkCherokeeFold(t, tc.pipe, tc.in, tc.want) })
	}
}

func TestCherokeeFoldPreservesOtherBytes(t *testing.T) {
	p := New().FoldCase()
	for b := range 256 {
		in := string(byte(b)) + "Ꮳ" + string(byte(b)) + "Ᏸ" + string(byte(b))
		want := cherokeeFoldOracle.Replace(cases.Fold().String(in))
		checkCherokeeFold(t, p, in, want)
	}
	// Adjacent code points must retain the dependency's behavior, too.
	for _, r := range []rune{0x139f, 0x13f6, 0x13f7, 0x13fe, 0xab6f, 0xabc0, utf8.RuneError} {
		in := string(r) + "Ꮳ" + string(r)
		want := cherokeeFoldOracle.Replace(cases.Fold().String(in))
		checkCherokeeFold(t, p, in, want)
	}
}

func TestFoldCherokeePrefilter(t *testing.T) {
	var latinExtendedE, uppercaseCherokee strings.Builder
	for r := '\uAB40'; r <= '\uAB6F'; r++ {
		latinExtendedE.WriteRune(r)
	}
	for r := '\u13C0'; r <= '\u13F5'; r++ {
		uppercaseCherokee.WriteRune(r)
	}
	latin := latinExtendedE.String()
	uppercase := uppercaseCherokee.String()

	for _, tc := range []struct {
		name     string
		in, want string
	}{
		{"no-candidate", "café 中文\xff\x80", "café 中文\xff\x80"},
		{"both-leads-without-prefix", "café \xe1!\xea?\x8f\xad\xae", "café \xe1!\xea?\x8f\xad\xae"},
		{"korean", "한국어", "한국어"},
		{"greek-extended", "ἀἄἆ", "ἀἄἆ"},
		{"korean-and-greek-extended", "한국어 ἀἄἆ\xff\x80", "한국어 ἀἄἆ\xff\x80"},
		{"latin-extended-e", latin, latin},
		{"latin-extended-e-with-cherokee", latin + "\uab70\uab80ᏸ" + latin, latin + "ᎠᎰᏰ" + latin},
		{"uppercase-cherokee", uppercase, uppercase},
		{"ea-ad-first", "café \uab70\uab80ᏸ", "café ᎠᎰᏰ"},
		{"ea-ae-before-ea-ad", "café \uab80\uab70ᏸ", "café ᎰᎠᏰ"},
		{"e1-before-ea-ad", "café ᏸ\uab70\uab80", "café ᏰᎠᎰ"},
		{"ea-first", "café ꮳᏸ", "café ᏣᏰ"},
		{"e1-first", "café ᏸꮳ", "café ᏰᏣ"},
		{"non-cherokee-ea", "café 가ᏸ", "café 가Ᏸ"},
		{"non-cherokee-e1", "café \u1200ꮳ", "café \u1200Ꮳ"},
		{"only-non-cherokee", "café 가\u1200", "café 가\u1200"},
		{"invalid-ea", "café \xeaꮳᏸ", "café \xeaᏣᏰ"},
		{"invalid-e1", "café \xe1ᏸꮳ", "café \xe1ᏰᏣ"},
		{"invalid-prefix", "café \xf0\x80ꮳ", "café \xf0\x80Ꮳ"},
		{"incomplete-ea-before-mapping", "café \xea\xadᏸ", "café \xea\xadᏰ"},
		{"incomplete-e1-before-mapping", "café \xe1\x8fꮳ", "café \xe1\x8fᏣ"},
		{"split-ea-ad", "café \xea\xff\xad\xb0", "café \xea\xff\xad\xb0"},
		{"split-ea-ad-before-mapping", "café \xea!\xad\xb0\uab70", "café \xea!\xad\xb0Ꭰ"},
		{"split-ea-ad-after-mapping", "café \uab70\xea\xff\xad\xb0", "café Ꭰ\xea\xff\xad\xb0"},
		{"truncated-ea", "café \xea", "café \xea"},
		{"truncated-ea-ad", "café \xea\xad", "café \xea\xad"},
		{"truncated-ea-ae", "café \xea\xae", "café \xea\xae"},
		{"truncated-e1", "café \xe1", "café \xe1"},
		{"truncated-e1-8f", "café \xe1\x8f", "café \xe1\x8f"},
		{"truncated-ea-ad-after-mapping", "café \uab70\xea\xad", "café Ꭰ\xea\xad"},
		{"truncated-after-mapping", "café ꮳ\xe1\x8f", "café Ꮳ\xe1\x8f"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := foldCherokee(tc.in); got != tc.want {
				t.Fatalf("foldCherokee(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFoldCaseCorrectionAllocations(t *testing.T) {
	p := New().FoldCase()
	for _, in := range []string{
		"", "ascii", "ASCII", "café 中文", "\xff\x80", "Straße Σς", "ꮳ", "Straße ꮳ", "Ꮳ", "Straße Ꮳ",
		"café 가\u1200", "café \xea", "café \xea\xad", "café \xea\xae", "café \xe1", "café \xe1\x8f",
		"한국어 ἀἄἆ", "\uab40\uab6f", "\uab40\uab70\uab80ᏸ\uab6f", "\u13c0\u13f5", "\xea\xff\xad\xb0",
		"café \xe1!\xea?\x8f\xad\xae",
	} {
		t.Run(fmt.Sprintf("%q", in), func(t *testing.T) {
			var sink string
			base := testing.AllocsPerRun(20, func() { sink = cases.Fold().String(in) })
			got := testing.AllocsPerRun(20, func() { sink, _ = p.Run(in) })
			allowed := base
			if folded := cases.Fold().String(in); cherokeeFoldOracle.Replace(folded) != folded {
				allowed++ // One correction buffer only when the folded output needs it.
			}
			if got > allowed {
				t.Errorf("allocations = %g, want <= %g (dependency = %g)", got, allowed, base)
			}
			t.Logf("dependency=%g corrected=%g allocs/run", base, got)
			runtime.KeepAlive(sink)
		})
	}
}
