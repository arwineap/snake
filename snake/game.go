package snake

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewGame(logger *zap.Logger) (*Game, error) {
	// random seed for apple gen
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		screenWidth:  960,
		screenHeight: 960,
		padding:      20,
		logger:       logger,

		Snake: Snake{
			Speed:     time.Millisecond * 100,
			Direction: right,
			Width:     16,
		},

		Apple: Apple{
			DropFrequency: time.Second,
			Width:         16,
		},

		Score: Score{
			Count: 0,
		},

		Border: Border{
			Padding: 20,
			Width:   2,
		},
	}

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Snake")

	// Reset State
	game.Restart()

	// Setup score
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return &Game{}, fmt.Errorf("could not setup font: %w", err)
	}
	game.Score.Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16.0,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return game, fmt.Errorf("could not configure font: %w", err)
	}

	return game, nil
}

type Game struct {
	screenWidth  int
	screenHeight int
	padding      int

	Snake Snake
	Apple Apple

	Score Score

	Border Border

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
	// Reset snake position
	g.Snake.Position = []Point{{(g.Snake.Width * 3) + g.padding, (g.padding + 1) + (g.Snake.Width / 2), g.Snake.Width}, {(g.Snake.Width * 2) + g.padding, (g.padding + 1) + (g.Snake.Width / 2), g.Snake.Width}}
	// Reset apples
	g.Apple.Positions = []Point{}
	// Reset Score
	g.Score.Count = 0
	// Reset direction
	g.Snake.Direction = right
	// Reset game
	g.over = false
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

	if g.checkCollisionWall() {
		g.gameOver("wall collision")
		return nil
	}

	if g.Snake.PointCollides(g.Snake.NextPoint()) {
		g.gameOver("snake collision")
		return nil
	}
	moved := g.Snake.Move()

	appleCollisions := g.Apple.Collisions(g.Snake.Head())
	for _, rmApple := range appleCollisions {
		g.Apple.Remove(rmApple)
		g.Score.Count++
		g.Snake.AddExtend(1)
	}

	// We only extend the snake if this is the snake has moved this tick
	if moved {
		g.Snake.Extend()
	}

	if g.Apple.CheckDrop() {
		g.Apple.Drop(g.randomUnusedPoint(g.Apple.Width))
	}

	return nil
}

func (g *Game) updateDirection() {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.Snake.Direction != left {
			g.Snake.SetDirection(right)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.Snake.Direction != up {
			g.Snake.SetDirection(down)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.Snake.Direction != right {
			g.Snake.SetDirection(left)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.Snake.Direction != down {
			g.Snake.SetDirection(up)
		}
	}
}

func (g *Game) randomUnusedPoint(width int) Point {
	randNumInsideMap := func() int {
		return rand.Intn((g.screenWidth-g.padding)-g.padding) + g.padding
	}

	var point = Point{X: randNumInsideMap(), Y: randNumInsideMap(), Width: width}
	for g.Snake.PointCollides(point) && !g.Apple.PointCollides(point) {
		point = Point{X: randNumInsideMap(), Y: randNumInsideMap(), Width: width}
	}

	return point
}

func (g *Game) gameOver(reason string) {
	g.logger.Info("game over", zap.Any("reason", reason), zap.Any("score", g.Score.Count))
	g.over = true
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	w, h := screen.Size()

	g.Score.Draw(screen)
	text.Draw(screen, "game over", g.Score.Font, int(float64(w)*0.46), int(float64(h)*0.5), color.RGBA{R: 255, G: 100, B: 100, A: 255})
	text.Draw(screen, "press r to restart", g.Score.Font, int(float64(w)*0.43), int(float64(h)*0.55), color.White)
	return
}

func (g *Game) checkCollisionWall() bool {
	// Head of the Snake is g.Snake.Position[0]
	// Tail of the Snake is g.Snake.Position[len(g.Snake.Position)-1]
	// check if head of Snake + Direction is within the "pixel" of death
	// We start in the top left (0, 0)
	return g.Border.Collides(g.Snake.Head())
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 0x00, G: 0x40, B: 0x80, A: 0xff})

	if g.over {
		g.drawGameOver(screen)
		return
	}

	// Draw border
	//borderLines := rect(float64(g.padding), float64(g.padding), float64(g.screenHeight-2*g.padding), float64(g.screenWidth-2*g.padding))
	// for _, line := range borderLines {
	// 	ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	// }
	g.Border.Draw(screen)

	g.Snake.Draw(screen)
	g.Apple.Draw(screen)
	g.Score.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}
