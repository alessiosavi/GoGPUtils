// Package textnorm provides fluent, deterministic text normalization pipelines.
//
// Use SplitTokens(), MapTokens(), FilterTokens(), and JoinTokens() when you need
// word-level cleanup. Use SearchPreset(), CanonicalPreset(), DBSafePreset(), and
// WithWidthFold() when you want ready-made pipelines instead of assembling stages
// by hand.
//
// Examples:
//
//	textnorm.SearchPreset().Run("  Café, go!  ")
//	textnorm.CanonicalPreset().Run("  Hello, World!  ")
//	textnorm.DBSafePreset(textnorm.WithWidthFold()).Run("  Ｇｏ\x00  ")
//
// These presets keep width folding opt-in and reuse the same stage API exposed
// by the base Pipeline type.
package textnorm
