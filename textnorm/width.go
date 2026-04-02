package textnorm

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

// FoldWidth appends an explicit width-folding stage.
func (p Pipeline) FoldWidth() Pipeline {
	return p.Then(func(s string) (string, error) {
		result, _, err := transform.String(width.Fold, s)
		if err != nil {
			return "", err
		}

		return result, nil
	})
}
