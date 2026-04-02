// Package textnorm provides fluent, deterministic text normalization pipelines.
//
// Use the core Pipeline and SplitTokens(), MapTokens(), FilterTokens(), and
// JoinTokens() helpers for explicit normalization flows. SearchPreset(),
// CanonicalPreset(), DBSafePreset(), and WithWidthFold() provide thin reusable
// presets for the common search, canonical, and persistence-safe cases.
//
// Benchmarks and fuzz targets live in this package so you can measure and harden
// the current implementation before any optimization work begins.
//
// Examples:
//
//	textnorm.SearchPreset().Run("  Café, go!  ")
//	textnorm.CanonicalPreset().Run("  Hello, World!  ")
//	textnorm.DBSafePreset(textnorm.WithWidthFold()).Run("  Ｇｏ\x00  ")
//
// Streaming adapters are intentionally deferred until real usage proves they are
// worth the extra surface area.
package textnorm
