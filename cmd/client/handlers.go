package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	defer fmt.Println("> ")
	return func(ps routing.PlayingState) {
		gs.HandlePause(ps)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	defer fmt.Println("> ")
	return func(move gamelogic.ArmyMove) {
		gs.HandleMove(move)
	}
}
