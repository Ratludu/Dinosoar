package game

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/ratludu/dinosoar/internal/entities"
)

type State int

const (
	Playing State = iota
	Paused
)

type Game struct {
	Settings       *Settings
	GameBackground *entities.Background
	Player         *entities.Player
	Pipes          []*entities.Pipes
	Count          int
	State          State
	FontSource     *text.GoTextFaceSource
	FontFace       *text.GoTextFace
}

func NewGame() *Game {

	game := &Game{
		Settings: NewSettings(),
	}

	var err error
	game.FontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}

	game.FontFace = &text.GoTextFace{
		Source: game.FontSource,
		Size:   16,
	}

	player, err := entities.NewPlayer(float64(game.Settings.screenWidth/5), float64(game.Settings.screenHeight/2))
	if err != nil {
		log.Fatal(err)
	}

	player.Score = 0

	game.Player = player

	background, err := entities.NewBackground()
	if err != nil {
		log.Fatal(err)
	}

	game.GameBackground = background

	pipeArray := make([]*entities.Pipes, 100)
	for i := 0; i < 100; i++ {
		pipes, err := entities.NewPipe()
		if err != nil {
			log.Fatal(err)
		}
		pipes.PipeX = float64(game.Settings.screenWidth) + float64(i*32*4)
		pipes.Height = float64(randBetween(game.Settings.screenHeight-int(pipes.Gap)-32, 32))
		pipeArray[i] = pipes
	}

	game.Pipes = pipeArray

	game.State = Paused

	return game
}

func (g *Game) Restart() {
	newGame := NewGame()

	g.Settings = newGame.Settings
	g.GameBackground = newGame.GameBackground
	g.Player = newGame.Player
	g.Pipes = newGame.Pipes
	g.Count = newGame.Count
	g.State = newGame.State
}

func (g *Game) Update() error {

	g.Count++
	if g.State == Paused {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = Playing
		}
		return nil
	}

	g.GameBackground.Update(g.Settings.screenWidth)
	score := 0
	for _, i := range g.Pipes {
		i.Update()
		if g.Player.X > i.PipeX {
			score++
		}

		playerRect := image.Rect(int(g.Player.X)+2, int(g.Player.Y)+2, int(g.Player.X)+24, int(g.Player.Y)+24)
		bottomPipe := image.Rect(int(i.PipeX)+8, int(float64(g.Settings.screenHeight)-24-i.Height), int(i.PipeX)+24, g.Settings.windowHeight)
		if collisions(playerRect, bottomPipe) {
			g.Restart()
		}

		topPipe := image.Rect(int(i.PipeX)+8, -1000, int(i.PipeX)+24, int(float64(g.Settings.screenHeight)-24-i.Height-i.Gap))
		if collisions(playerRect, topPipe) {
			g.Restart()
		}

	}

	g.Player.Score = score

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Player.Y += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Player.Y -= 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.Background(screen)
	for _, i := range g.Pipes {
		g.DrawPipes(i, screen)
	}
	g.AnimatePlayer(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("SCORE: %d", g.Player.Score))

	if g.State == Paused {
		pressSpace := "Press Space"
		toStart := "to Start"

		textWidth, textHeight := text.Measure(pressSpace, g.FontFace, 0)

		middle := textWidth / 2
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(float64(g.Settings.screenWidth)/2-middle, float64(g.Settings.screenHeight)/2-textHeight)
		text.Draw(screen, pressSpace, g.FontFace, textOp)

		textWidth, textHeight = text.Measure(toStart, g.FontFace, 0)
		middle = textWidth / 2
		textOp = &text.DrawOptions{}
		textOp.GeoM.Translate(float64(g.Settings.screenWidth)/2-middle, float64(g.Settings.screenHeight)/2+5)
		text.Draw(screen, toStart, g.FontFace, textOp)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.Settings.screenWidth), int(g.Settings.screenHeight)
}

func (g *Game) RunGame() {
	ebiten.SetWindowSize(int(g.Settings.windowWidth), int(g.Settings.windowHeight))
	ebiten.SetWindowTitle(g.Settings.Title)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) AnimatePlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(g.Settings.frameWidth)/2, -float64(g.Settings.frameHeight)/2)
	op.GeoM.Translate(float64(g.Settings.screenWidth/5), float64(g.Player.Y))
	i := (g.Count / 18) % g.Settings.frameCount
	sx, sy := g.Settings.frameOX+i*g.Settings.frameWidth, g.Settings.frameOY
	screen.DrawImage(g.Player.PlayerImage.SubImage(image.Rect(sx, sy, sx+g.Settings.frameWidth, sy+g.Settings.frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Background(screen *ebiten.Image) {

	bgWidth := g.GameBackground.BackgroundImage.Bounds().Dx()
	bgHeight := g.GameBackground.BackgroundImage.Bounds().Dy()

	scaleX := float64(g.Settings.screenWidth) / float64(bgWidth)
	scaleY := float64(g.Settings.screenHeight) / float64(bgHeight)
	scaledWidth := float64(bgWidth) * scaleX

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(-g.GameBackground.ScrollOffset, 0)
	screen.DrawImage(g.GameBackground.BackgroundImage, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(scaledWidth-g.GameBackground.ScrollOffset, 0)
	screen.DrawImage(g.GameBackground.BackgroundImage, op)
}

func (g *Game) DrawPipes(pipe *entities.Pipes, screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(pipe.PipeX), float64(g.Settings.screenHeight)-32-pipe.Height)
	screen.DrawImage(pipe.Pipe.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)

	for i := 0; i <= int(pipe.Height); i++ {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pipe.PipeX), float64(g.Settings.screenHeight)-32-pipe.Height+32+float64(i))
		screen.DrawImage(pipe.Pipe.SubImage(image.Rect(0, 31, 32, 32)).(*ebiten.Image), op)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(math.Pi)
	op.GeoM.Translate(float64(pipe.PipeX)+32, float64(g.Settings.screenHeight)-32-pipe.Height-pipe.Gap)
	screen.DrawImage(pipe.Pipe.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), op)

	reverseHeight := float64(g.Settings.screenHeight) - 32 - pipe.Height - pipe.Gap - 32
	for i := 0; i <= int(reverseHeight); i++ {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(math.Pi)
		op.GeoM.Translate(float64(pipe.PipeX)+32, float64(g.Settings.screenHeight)-32-pipe.Height-pipe.Gap-32-float64(i))
		screen.DrawImage(pipe.Pipe.SubImage(image.Rect(0, 31, 32, 32)).(*ebiten.Image), op)
	}
}
