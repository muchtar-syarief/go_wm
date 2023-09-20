package go_wm

import (
	"fmt"

	pfont "gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type WPosition string

const (
	TOP_LEFT      WPosition = "top_left"
	TOP_RIGHT     WPosition = "top_right"
	BOTTOM_LEFT   WPosition = "bottom_left"
	BOTTOM_RIGHT  WPosition = "bottom_right"
	CENTER        WPosition = "center"
	CENTER_LEFT   WPosition = "center_left"
	CENTER_RIGHT  WPosition = "center_right"
	CENTER_TOP    WPosition = "center_top"
	CENTER_BOTTOM WPosition = "center_bottom"
)

type WatermarkPosition struct {
	MinX pfont.Length
	MinY pfont.Length
	MaxX pfont.Length
	MaxY pfont.Length

	Length pfont.Length
	Height pfont.Length

	point vg.Point
}

func (w *WatermarkPosition) GetPosition(position WPosition) vg.Point {
	switch position {
	case TOP_LEFT:
		return w.topLeft()
	case TOP_RIGHT:
		return w.topRight()
	case BOTTOM_LEFT:
		return w.bottomLeft()
	case BOTTOM_RIGHT:
		return w.bottomRight()
	case CENTER:
		return w.center()
	case CENTER_LEFT:
		return w.centerLeft()
	case CENTER_RIGHT:
		return w.centerRight()
	case CENTER_TOP:
		return w.centerTop()
	case CENTER_BOTTOM:
		return w.centerBottom()
	}

	return w.center()
}

func (w *WatermarkPosition) topLeft() vg.Point {
	w.point = vg.Point{
		X: w.MinX,
		Y: w.MaxY - w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) bottomLeft() vg.Point {
	w.point = vg.Point{
		X: w.MinX,
		Y: w.MinY,
	}

	return w.point
}

func (w *WatermarkPosition) topRight() vg.Point {
	w.point = vg.Point{
		X: w.MaxX - w.Length,
		Y: w.MaxY - w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) bottomRight() vg.Point {
	w.point = vg.Point{
		X: w.MaxX - w.Length,
		Y: w.MinY,
	}

	return w.point
}

func (w *WatermarkPosition) center() vg.Point {
	w.point = vg.Point{
		X: w.MinX + (w.MaxX-w.MinX)/2 - w.Length/2,
		Y: w.MaxY/2 + w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) centerLeft() vg.Point {
	w.point = vg.Point{
		X: w.MinX,
		Y: w.MaxY/2 + w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) centerRight() vg.Point {
	w.point = vg.Point{
		X: w.MaxX - w.Length,
		Y: w.MaxY/2 + w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) centerTop() vg.Point {
	w.point = vg.Point{
		X: w.MinX + (w.MaxX-w.MinX)/2 - w.Length/2,
		Y: w.MaxY - w.Height,
	}

	return w.point
}

func (w *WatermarkPosition) centerBottom() vg.Point {
	fmt.Println(w.MinX)
	fmt.Println(w.MaxX)
	fmt.Println(w.Length)
	fmt.Println(w.MaxX - w.Length)
	fmt.Println((w.MaxX - w.Length) / 2)

	w.point = vg.Point{
		X: w.MinX + (w.MaxX-w.MinX)/2 - w.Length/2,
		Y: w.MinY,
	}

	return w.point
}
