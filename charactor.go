package stellar

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	CharactorDefaultHP     = 100
	CharactorDefaultAttack = 10
)

// charactor ...
type charactor struct {
	id    string
	name  string
	maxHP int64
	curHP int64
}

// CharactorName define a charactor name option for charactor.
func CharactorName(name string) CharactorOption {
	return func(c *charactor) error {
		c.name = name
		return nil
	}
}

// CharactorHP define a charactor hp option for charactor.
// The params can only deal with 1 or 2, otherwise returns error.
// If only 1 param, it will set both maxHP and curHP.
// If 2 param, hps[0] is maxHP, hps[1] is curHP.
func CharactorHP(hps ...int64) CharactorOption {
	return func(c *charactor) error {
		if len(hps) > 0 {
			c.maxHP = hps[0]
			c.curHP = hps[0]
			if len(hps) == 2 {
				c.curHP = hps[1]
			}
			return nil
		}
		return fmt.Errorf("the param number should be 1 or 2")
	}
}

// NewCharactor ...
func NewCharactor(opts ...CharactorOption) (*charactor, error) {
	id := uuid.New().String()
	c := &charactor{
		id:    id,
		name:  id[0:6],
		maxHP: CharactorDefaultHP,
		curHP: CharactorDefaultHP,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// GetName ...
func (c *charactor) GetName() string {
	return c.name
}

// GetMaxHP ...
func (c *charactor) GetMaxHP() int64 {
	return c.maxHP
}

// GetCurHP ...
func (c *charactor) GetCurHP() int64 {
	return c.curHP
}

// IsDead ...
func (c *charactor) IsDead() bool {
	return c.curHP <= 0
}

type attackSource struct {
	hp int64
}

func (c *charactor) doAttack() *attackSource {
	// TODO calc attack
	return &attackSource{hp: CharactorDefaultAttack}
}

type attackResult struct {
	expectHP int64
	actualHP int64
}

func (c *charactor) beAttack(src *attackSource) *attackResult {
	r := &attackResult{}
	r.expectHP = src.hp
	r.actualHP = min(c.curHP, r.expectHP)
	c.curHP -= r.actualHP
	return r
}
