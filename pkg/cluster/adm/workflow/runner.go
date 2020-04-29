package workflow

import (
	"ko3-gin/pkg/host"
)

type RunData = interface{}

type phaseRunner struct {
	Phase
	parent *phaseRunner
	level  int
	host   host.Host
}

func addPhaseRunner(r *Runner, parentRunner *phaseRunner, phase Phase) {
	currentRunner := &phaseRunner{
		Phase:  phase,
		parent: parentRunner,
		host:   r.host,
	}
	r.phaseRunners = append(r.phaseRunners, currentRunner)
	for _, childPhase := range phase.Phases {
		addPhaseRunner(r, currentRunner, childPhase)
	}
}

type Runner struct {
	Phases       []Phase
	phaseRunners []*phaseRunner
	host         host.Host
}

func NewRunner() *Runner {
	return &Runner{
		Phases: []Phase{},
	}
}

func (r *Runner) AppendRunner(phase Phase) {
	r.Phases = append(r.Phases, phase)
}

func (r *Runner) Run(data RunData, host host.Host) error {
	r.prepareForExecution()
	if err := r.visitAll(func(p *phaseRunner) error {
		if e := p.Run(data, host); e != nil {
			return e
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *Runner) prepareForExecution() {
	r.phaseRunners = []*phaseRunner{}
	var parentRunner *phaseRunner
	for _, phase := range r.Phases {
		addPhaseRunner(r, parentRunner, phase)
	}
}
func (r *Runner) visitAll(fn func(*phaseRunner) error) error {
	for _, currentRunner := range r.phaseRunners {
		if err := fn(currentRunner); err != nil {
			return err
		}
	}
	return nil
}
