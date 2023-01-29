package stellar

type attackAction struct {
	t       ActionType
	source  Charactor
	targets []Charactor
}

func AttackAction(source Charactor, targets ...Charactor) Action {
	return &attackAction{
		t:       ActionTypeAttack,
		source:  source,
		targets: targets,
	}
}

func (a *attackAction) Source() Charactor {
	return a.source
}

func (a *attackAction) Type() ActionType {
	return a.t
}

func (a *attackAction) Targets() []Charactor {
	return a.targets
}
