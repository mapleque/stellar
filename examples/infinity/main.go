package main

import (
	"fmt"
	"time"

	"github.com/mapleque/stellar"
)

const FrameDuration = time.Second
const ActionDuration = time.Second

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
		(&game).autoplay(ActionDuration)
	}()
	go func() {
		(&game).render(FrameDuration)
	}()
	if err := game.engine.Start(); err != nil {
		panic(err)
	}
}

func (g *Infinity) autoplay(d time.Duration) {
	g.actionInit()
	g.actionDo()
	g.actionClear()

	if g.player.IsDead() {
		if err := g.engine.Stop(); err != nil {
			panic(err)
		}
		return
	}

	time.AfterFunc(d, func() { g.autoplay(d) })
}

func (g *Infinity) actionInit() {
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

func (g *Infinity) actionClear() {
	g.events, g.lastEventId = g.engine.GetEvents(g.lastEventId)

	// remove dead enemys
	var enemys []stellar.Charactor
	for _, enemy := range g.enemys {
		if !enemy.IsDead() {
			enemys = append(enemys, enemy)
		} else {
			g.player.AddExp(30)
			if g.player.GetExp() >= 100 {
				overflow := g.player.GetExp() - 100
				g.player.LevelUp(1)
				if overflow > 0 {
					g.player.AddExp(overflow)
				}
			}
		}
	}
	if len(enemys) != len(g.enemys) {
		g.enemys = enemys
	}

	// clear actions
	g.actions = []stellar.Action{}
}

func (g *Infinity) actionDo() {
	g.engine.DoAction(g.actions...)
}

func (g *Infinity) randomEnemy() [][]stellar.CharactorOption {
	return [][]stellar.CharactorOption{{
		stellar.CharactorName("enemy"),
		stellar.CharactorHP(50),
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

func (g *Infinity) render(d time.Duration) {
	g.renderFrame()
	if g.player.IsDead() {
		return
	}
	time.AfterFunc(d, func() { g.render(d) })
}

func (g *Infinity) renderFrame() {
	// fmt.Printf(
	// 	"%+v HP: %+v/%+v\t",
	// 	g.player.GetName(),
	// 	g.player.GetCurHP(),
	// 	g.player.GetMaxHP(),
	// )
	// for _, enemy := range g.enemys {
	// 	fmt.Printf(
	// 		"\t%+v HP: %+v/%+v",
	// 		enemy.GetName(),
	// 		enemy.GetCurHP(),
	// 		enemy.GetMaxHP(),
	// 	)
	// }
	// fmt.Println()
	if len(g.events) > 0 {
		for _, event := range g.events {
			fmt.Printf("\t%s\n", event.Message())
		}
	}
}
