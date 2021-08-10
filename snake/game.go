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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewGame(logger *zap.Logger) (*Game, error) {
	game := &Game{
		screenWidth:  240,
		screenHeight: 240,
		padding:      20,
		logger:       logger,

		Snake: Snake{
			Width:     4,
			Speed:     time.Millisecond * 100,
			Direction: right,
		},

		Apple: Apple{
			DropFrequency: time.Second,
			Width:         4,
		},

		Score: Score{
			Count: 0,
		},
	}

	ebiten.SetWindowSize(game.screenWidth*2, game.screenHeight*2)
	ebiten.SetWindowTitle("Snake")

	// Reset State
	game.Restart()

	// Setup score
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return &Game{}, fmt.Errorf("could not setup font: %w", err)
	}
	game.Score.Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20.0,
		DPI:     40,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return game, fmt.Errorf("could not configure font: %w", err)
	}

	return game, nil
}

type Point struct {
	X int
	Y int
}

type direction string

const (
	right direction = "right"
	left  direction = "left"
	up    direction = "up"
	down  direction = "down"
)

type Game struct {
	screenWidth  int
	screenHeight int
	padding      int

	Snake Snake
	Apple Apple

	Score Score

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
	g.Snake.Position = []Point{{4 + g.padding, g.padding}, {g.padding, g.padding}}
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
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
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

	if g.checkSnakeHeadHitSnake() {
		g.gameOver("snake collision")
		return nil
	}

	g.moveSnake()

	if g.Apple.CheckDrop() {
		g.Apple.Drop(g.randomUnusedPoint())
	}

	return nil
}

// Checks if snake will run into itself
func (g *Game) checkSnakeHeadHitSnake() bool {
	var newPoint Point
	currentPoint := g.Snake.Position[0]
	switch g.Snake.Direction {
	case up:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y - 4
	case down:
		newPoint = currentPoint
		newPoint.Y = newPoint.Y + 4
	case left:
		newPoint = currentPoint
		newPoint.X = newPoint.X - 4
	case right:
		newPoint = currentPoint
		newPoint.X = newPoint.X + 4
	}

	if g.checkSnakePointCollide(newPoint) {
		return true
	} else {
		return false
	}
}

// Check if snake will hit an apple
func (g *Game) checkSnakeHeadHitApple() (bool, Point) {
	head := g.Snake.Position[0]
	for _, a := range g.Apple.Positions {
		if checkPointCollision(head, a) {
			return true, a
		}
	}

	return false, Point{}
}

func (g *Game) updateDirection() {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.Snake.Direction != left {
			g.Snake.Direction = right
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.Snake.Direction != up {
			g.Snake.Direction = down
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.Snake.Direction != right {
			g.Snake.Direction = left
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.Snake.Direction != down {
			g.Snake.Direction = up
		}
	}
}

func (g *Game) checkSnakePointCollide(point Point) bool {
	for _, p := range g.Snake.Position {
		if checkPointCollision(p, point) {
			return true
		}
	}
	return false
}

func (g *Game) checkApplePointCollide(point Point) bool {
	for _, p := range g.Apple.Positions {
		if checkPointCollision(p, point) {
			return true
		}
	}
	return false
}

func (g *Game) randomUnusedPoint() Point {
	// random seed for apple gen
	rand.Seed(time.Now().UnixNano())

	randNumInsideMap := func() int {
		return rand.Intn((g.screenWidth-g.padding)-g.padding) + g.padding
	}

	var point = Point{X: randNumInsideMap(), Y: randNumInsideMap()}
	for g.checkSnakePointCollide(point) && !g.checkApplePointCollide(point) {
		point = Point{X: randNumInsideMap(), Y: randNumInsideMap()}
	}

	return point
}

func (g *Game) moveSnake() {
	if time.Since(g.Snake.LastMove) < g.Snake.Speed {
		return
	}

	// To move the Snake we prepend a Point to the start of the slice, then remove the final Position of the slice
	// If we moved into an apple, we don't remove the final Position
	var newPoint Point
	head := g.Snake.Position[0]
	switch g.Snake.Direction {
	case up:
		newPoint = Point{X: head.X, Y: head.Y - g.Snake.Width}
	case down:
		newPoint = Point{X: head.X, Y: head.Y + g.Snake.Width}
	case left:
		newPoint = Point{X: head.X - g.Snake.Width, Y: head.Y}
	case right:
		newPoint = Point{X: head.X + g.Snake.Width, Y: head.Y}
	default:
		return
	}

	g.Snake.Position = append([]Point{newPoint}, g.Snake.Position...)

	if collided, applePosition := g.checkSnakeHeadHitApple(); collided {
		// Next slot is an apple, remove the apple
		g.Apple.Remove(applePosition)
		g.Score.Count++
	} else {
		// Next slot isn't an apple, so don't extend the snake
		g.Snake.Position = g.Snake.Position[:len(g.Snake.Position)-1]
	}
	g.Snake.LastMove = time.Now()
}

func (g *Game) gameOver(reason string) {
	g.logger.Info("game over", zap.Any("reason", reason), zap.Any("score", g.Score.Count))
	g.over = true
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.Score.Draw(screen)
	text.Draw(screen, "game over", g.Score.Font, 100, 80, color.RGBA{R: 255, G: 100, B: 100, A: 255})
	text.Draw(screen, "press r to restart", g.Score.Font, 80, 120, color.White)
	return
}

func (g *Game) checkCollisionWall() bool {
	// Head of the Snake is g.Snake.Position[0]
	// Tail of the Snake is g.Snake.Position[len(g.Snake.Position)-1]
	// check if head of Snake + Direction is within the "pixel" of death
	// We start in the top left (0, 0)
	head := g.Snake.Position[0]
	switch g.Snake.Direction {
	case up:
		if head.Y <= g.padding {
			return true
		}
	case down:
		if head.Y >= g.screenHeight-g.padding {
			return true
		}
	case left:
		if head.X <= g.padding {
			return true
		}
	case right:
		if head.X >= g.screenWidth-g.padding {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 0x00, G: 0x40, B: 0x80, A: 0xff})

	if g.over {
		g.drawGameOver(screen)
		return
	}

	// Draw border
	borderLines := rect(float64(g.padding), float64(g.padding), float64(g.screenHeight-2*g.padding), float64(g.screenWidth-2*g.padding))
	for _, line := range borderLines {
		ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	g.Snake.Draw(screen)
	g.Apple.Draw(screen)
	g.Score.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}
