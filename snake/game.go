package snake

import (
	"errors"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	ebiten.SetWindowTitle("snake")

	return game, nil
}

type Game struct {
	screenWidth  int
	screenHeight int
	padding      int

	pointerX int
	pointerY int

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

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.pointerX += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.pointerY += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.pointerX -= 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.pointerY -= 4
	}

	// +1/-1 is to stop player before it reaches the border
	if g.pointerX >= g.screenWidth-g.padding {
		g.pointerX = g.screenWidth - g.padding - 1
	}

	if g.pointerX <= g.padding {
		g.pointerX = g.padding + 1
	}

	if g.pointerY >= g.screenHeight-g.padding {
		g.pointerY = g.screenHeight - g.padding - 1
	}

	if g.pointerY <= g.padding {
		g.pointerY = g.padding + 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	// Draw player as a rect
	ebitenutil.DrawRect(screen, float64(g.pointerX)-2, float64(g.pointerY)-2, 4, 4, color.Black)
	ebitenutil.DrawRect(screen, float64(g.pointerX)-1, float64(g.pointerY)-1, 2, 2, color.RGBA{R: 100, G: 255, B: 100, A: 255})

	// Draw border
	borderLines := rect(float64(g.padding), float64(g.padding), float64(g.screenHeight-2*g.padding), float64(g.screenWidth-2*g.padding))
	for _, line := range borderLines {
		ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
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

// 4 pixel square for the player
