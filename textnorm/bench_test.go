package textnorm

import "testing"

func BenchmarkPipeline(b *testing.B) {
	cases := []struct {
		name  string
		pipe  Pipeline
		input string
	}{
		{name: "ascii-noop", pipe: New(), input: "plain ascii text"},
		{name: "mixed-unicode", pipe: New().NormalizeUnicode().FoldCase().TrimSpace().CollapseWhitespace(), input: "  Café   Straße   Ｇｏ  "},
	}

	for _, tc := range cases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = tc.pipe.Run(tc.input)
			}
		})
	}
}

func BenchmarkSearchPreset(b *testing.B) {
	inputs := []struct {
		name  string
		input string
	}{
		{name: "mixed-unicode", input: "  Café,   go!  gophers "},
	}

	for _, tc := range inputs {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = SearchPreset().Run(tc.input)
			}
		})
	}
}

func BenchmarkCanonicalPreset(b *testing.B) {
	inputs := []struct {
		name  string
		input string
	}{
		{name: "mixed-unicode", input: "  Hello,   World!  "},
	}

	for _, tc := range inputs {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = CanonicalPreset().Run(tc.input)
			}
		})
	}
}

func BenchmarkDBSafePreset(b *testing.B) {
	inputs := []struct {
		name  string
		input string
	}{
		{name: "dirty-bytes", input: string([]byte{'g', 'o', 0x00, 0xff, '!', 0xfe})},
	}

	for _, tc := range inputs {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = DBSafePreset().Run(tc.input)
			}
		})
	}
}
