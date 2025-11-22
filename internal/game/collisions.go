package game

import "image"

func collisions(r1, r2 image.Rectangle) bool {
	if r1.Max.X <= r2.Min.X || r1.Min.X >= r2.Max.X {
		return false
	}

	if r1.Max.Y <= r2.Min.Y || r1.Min.Y >= r2.Max.Y {
		return false
	}

	return true
}
