package go_wm

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strings"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

type Position struct {
	Position  WPosition
	Translate vg.Point
	Scale     vg.Length
}

type Watermark struct {
	Font     WFont
	Position *Position
	Color    *color.RGBA

	Text   string
	Rotate float64
}

// WaterMark for adding a watermark on the image
func (mark *Watermark) WaterMark(img image.Image, fr image.Image) (image.Image, error) {

	// image's length to canvas's length
	var bounds image.Rectangle
	if fr != nil {
		bounds = fr.Bounds()
	} else {
		bounds = img.Bounds()
	}
	w := vg.Length(bounds.Max.X) * vg.Inch / vgimg.DefaultDPI
	h := vg.Length(bounds.Max.Y) * vg.Inch / vgimg.DefaultDPI
	diagonal := vg.Length(math.Sqrt(float64(w*w + h*h)))

	// imgResize := resize.Resize(uint(w), uint(h), img, resize.Lanczos2)

	// create a canvas, which width and height are diagonal
	c := vgimg.New(diagonal, diagonal)

	minX := diagonal/2 - w/2
	maxX := diagonal/2 + w/2
	minY := diagonal/2 - h/2
	maxY := diagonal/2 + h/2

	// draw image on the center of canvas
	rect := vg.Rectangle{}
	rect.Min.X = minX
	rect.Min.Y = minY
	rect.Max.X = maxX
	rect.Max.Y = maxY
	c.DrawImage(rect, img)

	if fr != nil {
		c.DrawImage(rect, fr)
	}

	fnt, err := mark.Font.GetFont()
	if err != nil {
		return nil, err
	}

	// set the color of markText
	c.SetColor(mark.Color)

	p := WatermarkPosition{
		MinX:   minX,
		MinY:   minY,
		MaxX:   maxX,
		MaxY:   maxY,
		Length: fnt.Width(mark.Text),
		Height: fnt.Extents().Height,
	}

	c.Rotate(mark.Rotate)

	if mark.Text != "" {
		wmLoc := p.GetPosition(mark.Position.Position)
		if mark.Position.Translate != (vg.Point{}) {
			wmLoc = wmLoc.Add(mark.Position.Translate)
		}
		if mark.Position.Scale != 0 {
			wmLoc = wmLoc.Scale(mark.Position.Scale)
		}

		c.FillString(
			*fnt,
			wmLoc,
			mark.Text,
		)
	}

	// canvas writeto jpeg
	// canvas.img is private
	// so use a buffer to transfer
	jc := vgimg.PngCanvas{Canvas: c}
	buff := new(bytes.Buffer)
	_, err = jc.WriteTo(buff)
	if err != nil {
		return nil, err
	}

	img, _, err = image.Decode(buff)
	if err != nil {
		return nil, err
	}

	// get the center point of the image
	ctp := int(diagonal * vgimg.DefaultDPI / vg.Inch / 2)

	// // cutout the marked image
	size := bounds.Size()
	bounds = image.Rect(ctp-size.X/2, ctp-size.Y/2, ctp+size.X/2, ctp+size.Y/2)
	rv := image.NewRGBA(bounds)

	draw.Draw(rv, bounds, img, bounds.Min, draw.Src)

	return rv, nil
}

func (w *Watermark) MarkingPicture(filepath string, frame string) (image.Image, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	var fr image.Image
	if frame != "" {

		f, err := os.Open(frame)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		frImg, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}

		fr = frImg
	}

	img, err = w.WaterMark(img, fr)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// ext is extension ex: ".jpg", ".jpeg", and ".png"
func (w *Watermark) WriteToBuffer(img image.Image, ext string) (rv *bytes.Buffer, err error) {
	ext = strings.ToLower(ext)
	rv = new(bytes.Buffer)
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(rv, img, &jpeg.Options{Quality: 100})
	case ".png":
		err = png.Encode(rv, img)
	}
	return rv, err
}
