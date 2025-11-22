package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const filePipe = "./assets/images/PipeStyle1.png"

type Pipes struct {
	Pipe        *ebiten.Image
	PipeX       float64
	PipeY       float64
	ScrollSpeed float64
	Height      float64
	Gap         float64
	Passed      bool
}

func NewPipe() (*Pipes, error) {

	img, _, err := ebitenutil.NewImageFromFile(filePipe)
	if err != nil {
		return nil, err
	}

	return &Pipes{
		Pipe:        ebiten.NewImageFromImage(img),
		ScrollSpeed: 1.0,
		Gap:         64,
	}, nil
}

func (p *Pipes) Update() {
	p.PipeX -= p.ScrollSpeed
}
