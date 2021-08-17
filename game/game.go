package game

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"snake/apple"
	"snake/border"
	"snake/point"
	"snake/score"
	"snake/snake"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.uber.org/zap"
)

func NewGame(logger *zap.Logger) (*Game, error) {
	// random seed for apple gen
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		screenWidth:  960,
		screenHeight: 960,
		logger:       logger,
	}

	scor, err := score.New()
	if err != nil {
		return nil, fmt.Errorf("could not setup score: %w", err)
	}
	brdr := border.New()

	game.score = scor
	game.snake = snake.New(append(snake.DefaultOption(), snake.BorderPaddingOption(brdr.Padding))...)
	game.border = &brdr
	game.apple = apple.New()

	game.Modules = append(game.Modules, game.snake)
	game.Modules = append(game.Modules, game.apple)
	game.Modules = append(game.Modules, game.border)
	game.Modules = append(game.Modules, game.score)

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Snake")

	// Reset State
	game.Restart()

	return game, nil
}

type Game struct {
	screenWidth  int
	screenHeight int

	snake   *snake.Snake
	apple   *apple.Apple
	score   *score.Score
	border  *border.Border
	Modules []interface{}

	over bool

	logger *zap.Logger
}

func (g *Game) Run() error {
	g.logger.Info("run game")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}

func (g *Game) Restart() {
	for _, mod := range g.Modules {
		if r, ok := mod.(Reseter); ok {
			r.Reset()
		}
	}

	// Reset game
	g.over = false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 64, B: 128, A: 255})

	if g.over {
		g.drawGameOver(screen)
		return
	}

	for _, m := range g.Modules {
		if mod, ok := m.(Drawer); ok {
			mod.Draw(screen)
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.logger.Info("game ended by player pressing escape")
		return errors.New("game ended by player")
	}

	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.logger.Info("game restarted by player")
			g.Restart()
			return nil
		}
		return nil
	}

	g.updateDirection()

	for _, m := range g.Modules {
		if mod, ok := m.(Ender); ok {
			gameOver, gameOverString := mod.EndGame(g.snake.Head())
			if gameOver {
				g.gameOver(gameOverString)
				return nil
			}
		}
	}

	for _, m := range g.Modules {
		if mod, ok := m.(Extender); ok {
			extendAmount := mod.ExtendSnake(g.snake.Head())
			g.snake.AddExtend(extendAmount)
			g.score.Count = g.score.Count + extendAmount
		}
	}

	for _, m := range g.Modules {
		if mod, ok := m.(Droper); ok {
			// TODO
			mod.Drop(g.randomUnusedPoint(4))

		}
	}

	g.snake.Move()

	return nil
}

func (g *Game) updateDirection() {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.snake.Direction != snake.Left {
			g.snake.SetDirection(snake.Right)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.snake.Direction != snake.Up {
			g.snake.SetDirection(snake.Down)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.snake.Direction != snake.Right {
			g.snake.SetDirection(snake.Left)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.snake.Direction != snake.Down {
			g.snake.SetDirection(snake.Up)
		}
	}
}

func (g *Game) randomUnusedPoint(width int) point.Point {
	randNumInsideMap := func() int {
		return rand.Intn((g.screenWidth-g.border.Padding)-g.border.Padding) + g.border.Padding
	}

	var pnt = point.Point{X: randNumInsideMap(), Y: randNumInsideMap(), Width: width}
	for g.snake.PointCollides(pnt) && !g.apple.PointCollides(pnt) {
		pnt = point.Point{X: randNumInsideMap(), Y: randNumInsideMap(), Width: width}
	}

	return pnt
}

func (g *Game) gameOver(reason string) {
	g.logger.Info("game over", zap.Any("reason", reason), zap.Any("score", g.score.Count))
	g.over = true
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	w, h := screen.Size()

	g.score.Draw(screen)
	text.Draw(screen, "game over", g.score.Font, int(float64(w)*0.46), int(float64(h)*0.5), color.RGBA{R: 255, G: 100, B: 100, A: 255})
	text.Draw(screen, "press r to restart", g.score.Font, int(float64(w)*0.43), int(float64(h)*0.55), color.White)
	return
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}
