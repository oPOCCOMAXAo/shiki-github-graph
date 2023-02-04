package image

import "fmt"

type ColorHex struct {
	R, G, B byte
}

func (color ColorHex) String() string {
	return fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
}

func componentGradient(
	maxColor, gradientSize byte,
	value, maxValue int64,
) byte {
	return maxColor - byte(int64(gradientSize)*value/maxValue)
}

func Color(value, max int64) string {
	return ColorHex{
		R: componentGradient(248, 224, value, max),
		G: componentGradient(248, 96, value, max),
		B: componentGradient(248, 224, value, max),
	}.String()
}
