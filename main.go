package main

import (
	"fmt"

	"snake/snake"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("could not setup logger: %w", err))
	}

	game, err := snake.NewGame(logger)
	if err != nil {
		panic(fmt.Errorf("could not setup game: %w", err))
	}

	_ = game.Run()
}
