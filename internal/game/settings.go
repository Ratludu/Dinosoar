package game

import "math/rand/v2"

type Settings struct {
	Title        string
	windowWidth  int
	windowHeight int
	screenWidth  int
	screenHeight int
	frameWidth   int
	frameHeight  int
	frameCount   int
	frameOX      int
	frameOY      int
}

func NewSettings() *Settings {
	return &Settings{
		Title:        "Dinosoar",
		windowWidth:  640,
		windowHeight: 480,
		screenWidth:  320,
		screenHeight: 240,
		frameWidth:   32,
		frameHeight:  32,
		frameOX:      0,
		frameOY:      0,
		frameCount:   3,
	}
}

func randBetween(maximum, minimum int) int {
	return rand.IntN(maximum-minimum+1) + minimum
}
