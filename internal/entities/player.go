package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const playerFile = "./assets/images/FlyingCat.png"

type Player struct {
	PlayerImage *ebiten.Image
	Y           float64
	X           float64
	Score       int
}

func NewPlayer(X, Y float64) (*Player, error) {

	img, _, err := ebitenutil.NewImageFromFile(playerFile)
	if err != nil {
		return nil, err
	}

	return &Player{
		PlayerImage: ebiten.NewImageFromImage(img),
		Y:           Y,
		X:           X,
	}, nil
}
