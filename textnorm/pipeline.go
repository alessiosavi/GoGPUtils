package textnorm

// Pipeline holds an ordered list of normalization stages.
// The zero value is valid.
type Pipeline struct {
	stages []Stage
}

// New returns a new empty pipeline.
func New() Pipeline {
	return Pipeline{}
}

// Then returns a new pipeline with stage appended.
func (p Pipeline) Then(stage Stage) Pipeline {
	if stage == nil {
		return p
	}

	stages := make([]Stage, len(p.stages)+1)
	copy(stages, p.stages)
	stages[len(p.stages)] = stage

	return Pipeline{stages: stages}
}

// Run executes all stages in declaration order.
func (p Pipeline) Run(input string) (string, error) {
	result := input
	var err error

	for _, stage := range p.stages {
		result, err = stage(result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
