package stellar

import (
	"sync"
	"time"
)

type engine struct {
	sync.Mutex
	clockDuration time.Duration
	stop          bool

	charactors map[string]*charactor
	events     []Event
	actions    []Action
}

// New ...
func New(opts ...EngineOption) (Engine, error) {
	e := &engine{
		clockDuration: 500 * time.Millisecond,
		charactors:    map[string]*charactor{},
	}
	logger.Infof(
		"initial engine with default args:\n"+
			"\tclock duration: %s\n"+
			"......",
		e.clockDuration)
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// AddCharactor ...
func (e *engine) AddCharactor(opts ...CharactorOption) (Charactor, error) {
	e.Lock()
	defer e.Unlock()

	c, err := NewCharactor(opts...)
	if err != nil {
		return nil, err
	}
	e.charactors[c.id] = c
	e.addEvents(addCharactorEvent(c))
	logger.Infof("add charactor %s", c.id)
	return c, nil
}

// GetCharactor ...
func (e *engine) GetCharactor(id string) (Charactor, bool) {
	c, exist := e.charactors[id]
	return c, exist
}

// RemoveCharactor ...
func (e *engine) RemoveCharactor(id string) bool {
	e.Lock()
	defer e.Unlock()

	_, exist := e.charactors[id]
	logger.Infof("delete charactor %s %v", id, exist)
	if !exist {
		return exist
	}
	delete(e.charactors, id)
	return true
}

// DoAction ...
func (e *engine) DoAction(actions ...Action) error {
	e.Lock()
	defer e.Unlock()

	e.actions = append(e.actions, actions...)
	logger.Infof("recieve %v actions", len(actions))
	return nil
}

// GetEvents ...
func (e *engine) GetEvents(id int) (events []Event, lastEventId int) {
	return e.events[id:], len(e.events)
}

func (e *engine) addEvents(events ...Event) {
	e.events = append(e.events, events...)
	logger.Infof("add %v events", len(events))
}

// Start start the game
func (e *engine) Start() error {
	logger.Infof("clock start with duration: %s", e.clockDuration)
	for t := range time.NewTicker(e.clockDuration).C {
		logger.Debugf("clock ticker at %s", t)
		if e.stop {
			return nil
		}
		e.doAction()
	}
	return nil
}

// Stop stop the game
func (e *engine) Stop() error {
	e.stop = true
	return nil
}

func (e *engine) doAction() {
	e.Lock()
	defer e.Unlock()

	if len(e.actions) > 0 {
		logger.Infof("process %v actions", len(e.actions))
	}
	for _, action := range e.actions {
		switch action.Type() {
		case ActionTypeAttack:
			e.doAttackAction(action.(*attackAction))
		}
	}
	e.actions = []Action{}
}

func (e *engine) doAttackAction(action *attackAction) {
	src := action.source.(*charactor)
	if src.IsDead() {
		return
	}
	srcAttack := src.doAttack()
	e.addEvents(doAttackEvent(src, srcAttack))
	for _, tarI := range action.targets {
		tar := tarI.(*charactor)
		if tar.IsDead() {
			continue
		}
		tarAttack := tar.beAttack(srcAttack)
		e.addEvents(beAttackEvent(src, tar, srcAttack, tarAttack))
		if tar.IsDead() {
			e.addEvents(deadEvent(tar))
		}
	}
}
