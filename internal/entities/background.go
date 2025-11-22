package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const fileBackground = "./assets/images/Background.png"

type Background struct {
	BackgroundImage *ebiten.Image
	ScrollSpeed     float64
	ScrollOffset    float64
}

func NewBackground() (*Background, error) {

	img, _, err := ebitenutil.NewImageFromFile(fileBackground)
	if err != nil {
		return nil, err
	}

	return &Background{
		BackgroundImage: ebiten.NewImageFromImage(img),
		ScrollSpeed:     1.0,
		ScrollOffset:    0,
	}, nil

}

func (b *Background) Update(screenWidth int) {
	b.ScrollOffset += b.ScrollSpeed

	bgWidth := b.BackgroundImage.Bounds().Dx()
	scaleX := float64(screenWidth) / float64(bgWidth)
	scaledWidth := float64(bgWidth) * scaleX

	if b.ScrollOffset >= scaledWidth {
		b.ScrollOffset = 0
	}
}
