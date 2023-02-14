package stellar

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

const (
	CharactorInitHP         = 100
	CharactorMinDamage      = 10
	CharactorMaxDamage      = 20
	CharactorCrit           = 128
	CharactorLevel          = 1
	CharactorStrength       = 1
	CharactorStrengthGrowth = 1
	CharactorAgility        = 1
	CharactorAgilityGrowth  = 1

	Strength2HPCoefficient    = 10
	Agility2DamageCoefficient = 10
)

// charactor ...
type charactor struct {
	sync.Mutex

	id   string
	name string

	baseHP int64
	curHP  int64

	minDamage int64
	maxDamage int64

	// crit percent
	crit uint8

	exp      int64
	level    int
	maxLevel int

	prop       *charactorProp
	propGrowth *charactorProp
}

type charactorProp struct {
	// add maxHP
	strength               int64
	strength2HPCoefficient int64
	// add damage
	agility                   int64
	agility2DamageCoefficient int64
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
			c.baseHP = hps[0]
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
		id:     id,
		name:   id[0:6],
		baseHP: CharactorInitHP,
		curHP:  CharactorInitHP,

		minDamage: CharactorMinDamage,
		maxDamage: CharactorMaxDamage,

		crit: CharactorCrit,

		exp:      0,
		level:    1,
		maxLevel: CharactorLevel,

		prop: &charactorProp{
			strength:                  CharactorStrength,
			strength2HPCoefficient:    Strength2HPCoefficient,
			agility:                   CharactorAgility,
			agility2DamageCoefficient: Agility2DamageCoefficient,
		},
		propGrowth: &charactorProp{
			strength: CharactorStrengthGrowth,
			agility:  CharactorAgilityGrowth,
		},
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
	return c.maxHP()
}

// GetCurHP ...
func (c *charactor) GetCurHP() int64 {
	return c.curHP
}

// IsDead ...
func (c *charactor) IsDead() bool {
	return c.curHP <= 0
}

// Level ...
func (c *charactor) Level() int {
	return c.level
}

// LevelUp ...
func (c *charactor) LevelUp(level int) {
	c.levelUp(level)
}

// AddExp ...
func (c *charactor) AddExp(exp int64) {
	c.addExp(exp)
}

// GetExp ...
func (c *charactor) GetExp() int64 {
	return c.exp
}

func (c *charactor) levelUp(level int) {
	c.Lock()
	defer c.Unlock()

	c.level += level
	c.level = max(0, c.level)
	c.level = min(c.level, c.maxLevel)

	c.curHP = c.maxHP()
	c.exp = 0
}

func (c *charactor) addExp(exp int64) {
	c.Lock()
	defer c.Unlock()

	c.exp += exp
}

func (c *charactor) maxHP() int64 {
	return c.baseHP + c.prop.strength*c.prop.strength2HPCoefficient
}

type attackSource struct {
	hp   int64
	crit bool
}

func (c *charactor) doAttack() *attackSource {
	baseDamage := c.minDamage + rand.Int63n(c.maxDamage-c.minDamage)
	floatDamage := c.prop.agility * c.prop.agility2DamageCoefficient
	damage := baseDamage + floatDamage
	crit := rand.Intn(255) < int(c.crit)
	if crit {
		damage *= 2
	}
	return &attackSource{hp: damage, crit: crit}
}

type attackResult struct {
	expectHP int64
	actualHP int64
}

func (c *charactor) beAttack(src *attackSource) *attackResult {
	c.Lock()
	defer c.Unlock()

	r := &attackResult{}
	r.expectHP = src.hp
	r.actualHP = min(c.curHP, r.expectHP)
	c.curHP -= r.actualHP
	return r
}
