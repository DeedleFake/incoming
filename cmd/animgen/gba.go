package main

import (
	"image/color"
)

type GBAColor uint16

func (c GBAColor) RGBA() (r, g, b, a uint32) {
	a = 0xFFFF

	r = (uint32(c) & 31) * a / 31
	g = ((uint32(c) >> 5) & 31) * a / 31
	b = ((uint32(c) >> 10) & 31) * a / 31

	return
}

var GBAModel = color.ModelFunc(func(c color.Color) color.Color {
	r, g, b, a := c.RGBA()

	return GBAColor(((b * 31 / a) << 10) | ((g * 31 / a) << 5) | (r * 31 / a))
})
