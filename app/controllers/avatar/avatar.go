package avatar

import (
	"fmt"
	"image"
	"image/color/palette"
	"math/rand"
	"time"

	"github.com/issue9/identicon"
)

const AVATAR_SIZE = 290

// RandomImage generates and returns a random avatar image unique to input data
// in custom size (height and width).
func RandomImageSize(size int, data []byte) (image.Image, error) {
	randExtent := len(palette.WebSafe) - 32
	rand.Seed(time.Now().UnixNano())
	colorIndex := rand.Intn(randExtent)
	backColorIndex := colorIndex - 1
	if backColorIndex < 0 {
		backColorIndex = randExtent - 1
	}

	// Define size, background, and forecolor
	imgMaker, err := identicon.New(size,
		palette.WebSafe[backColorIndex], palette.WebSafe[colorIndex:colorIndex+32]...)
	if err != nil {
		return nil, fmt.Errorf("identicon.New: %v", err)
	}
	return imgMaker.Make(data), nil
}

// RandomImage generates and returns a random avatar image unique to input data
// in default size (height and width).
func RandomImage(data []byte) (image.Image, error) {
	return RandomImageSize(AVATAR_SIZE, data)
}
