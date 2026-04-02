package textnorm

import "testing"

func TestPipelineThenDoesNotMutateOriginal(t *testing.T) {
	base := New()
	derived := base.Then(func(s string) (string, error) {
		return s + "!", nil
	})

	gotBase, err := base.Run("go")
	if err != nil {
		t.Fatalf("base Run() error = %v", err)
	}
	if gotBase != "go" {
		t.Fatalf("base Run() = %q, want %q", gotBase, "go")
	}

	gotDerived, err := derived.Run("go")
	if err != nil {
		t.Fatalf("derived Run() error = %v", err)
	}
	if gotDerived != "go!" {
		t.Fatalf("derived Run() = %q, want %q", gotDerived, "go!")
	}
}

func TestPipelineRunsStagesInOrder(t *testing.T) {
	pipe := New().Then(func(s string) (string, error) {
		return s + "a", nil
	}).Then(func(s string) (string, error) {
		return s + "b", nil
	})

	got, err := pipe.Run("")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got != "ab" {
		t.Fatalf("Run() = %q, want %q", got, "ab")
	}
}
