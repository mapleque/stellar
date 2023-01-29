package main

import (
	"fmt"
	"time"

	"github.com/mapleque/stellar"
)

const FrameDuration = time.Second

type Infinity struct {
	engine stellar.Engine
	player stellar.Charactor

	enemysOpts [][]stellar.CharactorOption
	actions    []stellar.Action

	enemys      []stellar.Charactor
	events      []stellar.Event
	lastEventId int
}

func main() {
	stellar.UseStdLogger()
	stellar.SetLogLevel(stellar.Info, false)

	game := Infinity{}
	engine, err := stellar.New()
	if err != nil {
		panic(err)
	}
	game.engine = engine

	c, err := game.engine.AddCharactor()
	if err != nil {
		panic(err)
	}
	game.player = c
	go func() {
		(&game).frame(FrameDuration)
	}()
	if err := game.engine.Start(); err != nil {
		panic(err)
	}
}

func (g *Infinity) frame(d time.Duration) {
	fmt.Printf("---- new frame ----\n")
	g.frameInit()
	g.frameDisplay()
	g.frameAction()
	g.frameClear()
	g.frameDisplay()
	fmt.Printf("========\n")

	if g.player.IsDead() {
		fmt.Printf("Your are dead!\n")
		if err := g.engine.Stop(); err != nil {
			panic(err)
		}
		return
	}

	time.AfterFunc(d, func() { g.frame(d) })
}

func (g *Infinity) frameInit() {
	if len(g.enemys) == 0 {
		g.enemysOpts = g.randomEnemy()
		for _, opts := range g.enemysOpts {
			enemy, err := g.engine.AddCharactor(opts...)
			if err != nil {
				fmt.Printf("[Error] generat enemy error: %+v\n", err)
			} else {
				g.enemys = append(g.enemys, enemy)
			}
		}
	}
	g.actions = g.randomActions()
	g.events = []stellar.Event{}
}

func (g *Infinity) frameClear() {
	g.events, g.lastEventId = g.engine.GetEvents(g.lastEventId)

	// remove dead enemys
	var enemys []stellar.Charactor
	for _, enemy := range g.enemys {
		if !enemy.IsDead() {
			enemys = append(enemys, enemy)
		}
	}
	if len(enemys) != len(g.enemys) {
		g.enemys = enemys
	}

	// clear actions
	g.actions = []stellar.Action{}
}

func (g *Infinity) frameDisplay() {
	fmt.Printf(
		"Player:\n\t%+v HP: %+v/%+v\n",
		g.player.GetName(),
		g.player.GetCurHP(),
		g.player.GetMaxHP(),
	)
	fmt.Printf("Enemy(s) (%d):\n", len(g.enemys))
	for _, enemy := range g.enemys {
		fmt.Printf(
			"\t%+v HP: %+v/%+v\n",
			enemy.GetName(),
			enemy.GetCurHP(),
			enemy.GetMaxHP(),
		)
	}
	if len(g.actions) > 0 {
		fmt.Printf("Actions:\n")
		for _, action := range g.actions {
			fmt.Printf(
				"\t%+v %+v %+v\n",
				action.Source(),
				action.Type(),
				action.Targets(),
			)
		}
	}
	if len(g.events) > 0 {
		fmt.Printf("Events:\n")
		for _, event := range g.events {
			fmt.Printf("\t%s\n", event.Message())
		}
	}
}

func (g *Infinity) frameAction() {
	g.engine.DoAction(g.actions...)
}

func (g *Infinity) randomEnemy() [][]stellar.CharactorOption {
	return [][]stellar.CharactorOption{{
		stellar.CharactorName("enemy"),
		stellar.CharactorHP(10),
	}}
}

func (g *Infinity) randomActions() []stellar.Action {
	actions := []stellar.Action{
		stellar.AttackAction(g.player, g.enemys...),
	}
	for _, enemy := range g.enemys {
		actions = append(actions, stellar.AttackAction(enemy, g.player))
	}
	return actions
}
