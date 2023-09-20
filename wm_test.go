package go_wm_test

import (
	"image/color"
	"os"
	"path"
	"strings"
	"testing"

	wm "github.com/muchtar-syarief/go_wm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
)

func TestWatermark(t *testing.T) {
	target := "./assets"
	fname := "roti.jpg"
	frame := "frame/Picture16.png"
	outputLoc := "./results"

	w := wm.Watermark{
		Font: wm.WFont{
			// Uri:      path.Join(target, "font/Times New Normal Regular.ttf"),
			// TypeFace: "Times New Normal",
			Size:   72,
			Weight: font.WeightExtraBold,
			Style:  font.StyleNormal,
		},
		Text: "ROTI JUICY",
		Position: &wm.Position{
			Position: wm.CENTER,
		},
		Color: &color.RGBA{0, 0, 0, 255},
		// Rotate: math.Pi * rand.Float64() / 2,
	}

	img, err := w.MarkingPicture(path.Join(target, fname), path.Join(target, frame))
	assert.Nil(t, err)

	ext := path.Ext(fname)
	buff, err := w.WriteToBuffer(img, ext)
	assert.Nil(t, err)

	base := strings.Split(fname, ".")[0] + "_marked"
	output := path.Join(outputLoc, base+ext)
	f, err := os.Create(output)
	assert.Nil(t, err)

	_, err = buff.WriteTo(f)
	assert.Nil(t, err)
}
