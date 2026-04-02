package textnorm

import (
	"strings"
)

// TokenStage transforms token slices and may return an error.
type TokenStage func([]string) ([]string, error)

// TokenPipeline holds ordered token stages derived from a string pipeline.
type TokenPipeline struct {
	source Pipeline
	stages []TokenStage
}

// SplitTokens turns the current string pipeline into a token pipeline.
func (p Pipeline) SplitTokens() TokenPipeline {
	return TokenPipeline{source: p}
}

// Then returns a new token pipeline with a stage appended.
func (tp TokenPipeline) Then(stage TokenStage) TokenPipeline {
	if stage == nil {
		return tp
	}

	stages := make([]TokenStage, len(tp.stages)+1)
	copy(stages, tp.stages)
	stages[len(tp.stages)] = stage

	return TokenPipeline{source: tp.source, stages: stages}
}

// MapTokens returns a new token pipeline that maps every token.
func (tp TokenPipeline) MapTokens(fn func(string) string) TokenPipeline {
	if fn == nil {
		return tp
	}

	return tp.Then(func(tokens []string) ([]string, error) {
		out := make([]string, len(tokens))
		for i, token := range tokens {
			out[i] = fn(token)
		}
		return out, nil
	})
}

// FilterTokens returns a new token pipeline that keeps matching tokens.
func (tp TokenPipeline) FilterTokens(fn func(string) bool) TokenPipeline {
	if fn == nil {
		return tp
	}

	return tp.Then(func(tokens []string) ([]string, error) {
		out := make([]string, 0, len(tokens))
		for _, token := range tokens {
			if fn(token) {
				out = append(out, token)
			}
		}
		return out, nil
	})
}

// Run executes the source pipeline, tokenizes it, and applies token stages.
func (tp TokenPipeline) Run(input string) ([]string, error) {
	text, err := tp.source.Run(input)
	if err != nil {
		return nil, err
	}

	tokens := strings.Fields(text)
	for _, stage := range tp.stages {
		tokens, err = stage(tokens)
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}

// JoinTokens joins token output back into a string pipeline.
func (tp TokenPipeline) JoinTokens(sep string) Pipeline {
	return tp.source.Then(func(input string) (string, error) {
		tokens, err := tp.Run(input)
		if err != nil {
			return "", err
		}

		return strings.Join(tokens, sep), nil
	})
}
