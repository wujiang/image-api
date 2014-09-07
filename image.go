package main

import (
	"image"
	"image/color"

	"github.com/golang/glog"
	"github.com/wujiang/imaging"
)

type Actions map[string]func(image.Image, int, int, imaging.ResampleFilter) *image.NRGBA

const (
	METHOD_WIDTH  = "width"  // use width and keep aspect ratio
	METHOD_SQUARE = "square" // square and crop
	METHOD_AUTO   = "auto"   // use width and get a square
)

var (
	squaredMethods = []string{
		METHOD_SQUARE,
		METHOD_AUTO,
	}

	actions = Actions{
		METHOD_WIDTH:  imaging.Resize,
		METHOD_SQUARE: imaging.Thumbnail,
		METHOD_AUTO:   imaging.Fit,
	}
)

func Resize(t string, src string, dst string, w int, h int, ot int64) (err error) {
	Action := actions[t]
	if Action == nil {
		Action = imaging.Resize
	}
	img, err := imaging.Open(src)
	if err != nil {
		glog.Warning(err)
		return
	}

	// exif standard
	// http://www.daveperrett.com/articles/2012/07/28/exif-orientation-handling-is-a-ghetto/
	switch {
	case ot == 2:
		img = imaging.FlipH(img)
	case ot == 3:
		img = imaging.Rotate180(img)
	case ot == 4:
		img = imaging.FlipH(img)
		img = imaging.Rotate180(img)
	case ot == 5:
		img = imaging.FlipV(img)
		img = imaging.Rotate270(img)
	case ot == 6:
		img = imaging.Rotate270(img)
	case ot == 7:
		img = imaging.FlipV(img)
		img = imaging.Rotate90(img)
	case ot == 8:
		img = imaging.Rotate90(img)
	}

	thumb := Action(img, w, h, imaging.Lanczos)
	// thumb := Action(img, w, h, imaging.CatmullRom)
	thumb = imaging.Sharpen(thumb, cfg.SharpenSigma)

	// get the width and height
	rect := thumb.Rect
	w = rect.Dx()
	h = rect.Dy()
	// create a new blank image
	newImg := imaging.New(w, h, color.NRGBA{0, 0, 0, 0})
	// paste the NARGB image into the background
	newImg = imaging.Paste(newImg, thumb, image.Pt(0, 0))
	if err = imaging.Save(newImg, dst, cfg.JpegQuality); err != nil {
		glog.Warning(err)
	}
	return
}
