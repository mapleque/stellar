package stellar

import "fmt"

type event struct {
	t          EventType
	msg        string
	charactors []Charactor
}

func (e *event) Type() EventType {
	return e.t
}

func (e *event) Charactors() []Charactor {
	return e.charactors
}

func (e *event) Message() string {
	return e.msg
}

func addCharactorEvent(c *charactor) Event {
	return &event{
		t:          EventTypeAddCharactor,
		msg:        fmt.Sprintf("add charactor %s", c.name),
		charactors: []Charactor{c},
	}
}

func doAttackEvent(c *charactor, attack *attackSource) Event {
	msg := fmt.Sprintf("%s do attack with damage %+v", c.name, attack.hp)
	if attack.crit {
		msg += "(crit)"
	}
	return &event{
		t:          EventTypeDoAttack,
		msg:        msg,
		charactors: []Charactor{c},
	}
}

func beAttackEvent(
	src, tar *charactor, sa *attackSource, ta *attackResult,
) Event {
	return &event{
		t: EventTypeBeAttack,
		msg: fmt.Sprintf(
			"%s demaged from %s %+v (overflow %v)",
			src.name,
			tar.name,
			sa.hp,
			ta.expectHP-ta.actualHP,
		),
		charactors: []Charactor{src, tar},
	}
}

func deadEvent(tar *charactor) Event {
	return &event{
		t: EventTypeDead,
		msg: fmt.Sprintf(
			"%s dead",
			tar.name,
		),
		charactors: []Charactor{tar},
	}
}
