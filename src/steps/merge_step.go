package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// MergeStep merges the branch with the given name into the current branch.
type MergeStep struct {
	Branch      domain.BranchName
	previousSha domain.SHA
	EmptyStep
}

func (step *MergeStep) CreateAbortSteps() []Step {
	return []Step{&AbortMergeStep{}}
}

func (step *MergeStep) CreateContinueSteps() []Step {
	return []Step{&ContinueMergeStep{}}
}

func (step *MergeStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToShaStep{Hard: true, Sha: step.previousSha}}, nil
}

func (step *MergeStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSha, err = run.Backend.CurrentSha()
	if err != nil {
		return err
	}
	return run.Frontend.MergeBranchNoEdit(step.Branch)
}
