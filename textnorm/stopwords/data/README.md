# Stopword data

The files `english`, `french`, and `italian` are taken verbatim from the
[NLTK stopwords corpus](https://github.com/nltk/nltk_data) and embedded
into the Go binary via `//go:embed` in `../stopwords.go`.

## Refresh from upstream

```bash
wget https://github.com/nltk/nltk_data/raw/refs/heads/gh-pages/packages/corpora/stopwords.zip
unzip -o stopwords.zip
cp stopwords/{english,french,italian} textnorm/stopwords/data/
```

After refreshing, run `go test ./textnorm/stopwords/...` from the repo root
to confirm the anchor-word tests still pass.

## Format

One word per line. Lines beginning with `#` and blank lines are ignored
by the parser. Surrounding whitespace is trimmed.

## Adding a new language

1. Drop the new file under this directory (e.g. `data/spanish`).
2. Add it to the `//go:embed` directive in `../stopwords.go`.
3. Add an accessor function (e.g. `func Spanish() map[string]struct{}`)
   following the lazy-`sync.Once` pattern used for the existing languages.
4. Add a `case` to `CleanAllStopwords` in `../stopwords.go`.
5. Add anchor-word and identity tests in `../stopwords_test.go`.
