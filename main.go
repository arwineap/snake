package main

import (
	"fmt"
	"snake/game"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("could not setup logger: %w", err))
	}

	g, err := game.NewGame(logger)
	if err != nil {
		panic(fmt.Errorf("could not setup game: %w", err))
	}

	_ = g.Run()
}
