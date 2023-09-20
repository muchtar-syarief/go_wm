# Add Watermark and Frame To Image

### Usage 
This Go packga use for adding Text watermark to an image or adding frame to an image.<br>
Size image output is following size of frame.<br>
<br>
For using :<br>
```
func main() {
    target := "./assets"
	fname := "roti.jpg"
	frame := "frame/Picture16.png"

    w := wm.Watermark{
        // Custom your Font for watermark
		Font: wm.WFont{
			Uri:      path.Join(target, "font/Times New Normal Regular.ttf"),
			TypeFace: "Times New Normal",
			Size:     72,
			Weight:   font.WeightExtraBold,
			Style:    font.StyleNormal,
		},
        // Text of watermark
		Text: "ROTI JUICY",
		Position: &wm.Position{
			Position:  wm.CENTER,
			Translate: vg.Point{},
		},
        // color font
		Color: &color.RGBA{0, 0, 0, 255},
		// Rotate: math.Pi * rand.Float64() / 2,
	}

    // Processing add watermark and frame
	img, err := w.MarkingPicture(path.Join(target, fname), path.Join(target, frame))
	assert.Nil(t, err)

    // geting buffer img 
	ext := path.Ext(fname)
	buff, err := w.WriteToBuffer(img, ext)
	assert.Nil(t, err)

    // save image
	base := strings.Split(fname, ".")[0] + "_marked"
	f, err := os.Create(base + ext)
	assert.Nil(t, err)

	_, err = buff.WriteTo(f)
	assert.Nil(t, err)
}
```


## Features
- [x] Custom Font
- [x] Custom Watermark location. ex: TOP_LEFT, CENTER, BOTTOM_RIGHT
- [x] Translte Text location.
- [x] Support ".jpg", ".jpeg", and ".png" extension
- [ ] Rotate Text
- [ ] Font Weight