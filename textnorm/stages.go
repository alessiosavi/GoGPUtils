package textnorm

// Stage transforms input text and may return an error.
type Stage func(string) (string, error)
