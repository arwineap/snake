package snake

import (
	"errors"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"go.uber.org/zap"
)

func NewGame(logger *zap.Logger) (*Game, error) {
	game := &Game{
		screenWidth:  240,
		screenHeight: 240,
		padding:      20,
		logger:       logger,
	}

	ebiten.SetWindowSize(game.screenWidth*2, game.screenHeight*2)
	ebiten.SetWindowTitle("Snake")
	// TODO defaults to 60
	ebiten.SetMaxTPS(10)

	// Setup initial Snake
	game.Snake.Position = append(game.Snake.Position, Point{4 + game.padding, game.padding}, Point{game.padding, game.padding})
	game.Snake.Width = 4

	//game.Snake.Direction = right

	// Setup initial apples
	game.Apple.DropFrequency = time.Second
	game.Apple.Width = 4

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

	logger *zap.Logger
}

func (g *Game) Run() error {
	g.logger.Info("run game")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.logger.Info("game ended by player pressing escape")
		return errors.New("game ended by player")
	}

	g.updateDirection()

	if g.checkCollisionWall() {
		g.gameOver("wall collision")
	}

	if g.checkSnakeHeadHitSnake() {
		g.gameOver("snake collision")
	}

	g.moveSnake()

	if g.Apple.CheckDrop() {
		g.Apple.Drop(g.randomUnusedPoint())
	}

	// g.logger.Debug("update debugger", zap.Reflect("Snake", g.Snake.Position))
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
	} else {
		// Next slot isn't an apple, so don't extend the snake
		g.Snake.Position = g.Snake.Position[:len(g.Snake.Position)-1]
	}
}

func (g *Game) gameOver(reason string) {
	g.logger.Fatal("game over", zap.String("reason", reason), zap.Any("score", len(g.Snake.Position)))
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

	// Draw border
	borderLines := rect(float64(g.padding), float64(g.padding), float64(g.screenHeight-2*g.padding), float64(g.screenWidth-2*g.padding))
	for _, line := range borderLines {
		ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	// Draw player
	g.drawSnake(screen)

	g.drawApples(screen)

	// Need a way to draw the whole Snake
}

func (g *Game) drawApples(screen *ebiten.Image) {
	for _, p := range g.Apple.Positions {
		ebitenutil.DrawRect(screen, float64(p.X)-2, float64(p.Y)-2, float64(g.Apple.Width), float64(g.Apple.Width), color.Black)
		ebitenutil.DrawRect(screen, float64(p.X)-1, float64(p.Y)-1, float64(g.Apple.Width)/2, float64(g.Apple.Width)/2, color.RGBA{R: 255, G: 100, B: 100, A: 255})
	}
}

func (g *Game) drawSnake(screen *ebiten.Image) {
	for _, p := range g.Snake.Position {
		ebitenutil.DrawRect(screen, float64(p.X)-2, float64(p.Y)-2, float64(g.Snake.Width), float64(g.Snake.Width), color.Black)
		ebitenutil.DrawRect(screen, float64(p.X)-1, float64(p.Y)-1, float64(g.Snake.Width)/2, float64(g.Snake.Width)/2, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}

type line struct {
	X1, Y1, X2, Y2 float64
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}
