package main

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/rwcarlsen/goexif/exif"
)

var (
	maxMemory = int64(10 << 20) // Byte, 10MB
)

func ParseResizerParams(r *http.Request) (method string, width, height int,
	file multipart.File, err error) {
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		return
	}

	m := r.MultipartForm
	methods := m.Value["method"]
	if len(methods) != 1 {
		err = errors.New("Unknown method")
		return
	}
	method = methods[0]
	widths := m.Value["width"]
	if len(widths) != 1 {
		err = errors.New("Unknown width")
		return
	}
	width, err = strconv.Atoi(widths[0])
	if err != nil {
		return
	}
	// time to figure out height
	if method == METHOD_SQUARE || method == METHOD_AUTO {
		height = width
	} else {
		height = 0
	}

	files := m.File["file"]
	if len(files) != 1 {
		err = errors.New("Expect one file")
		return
	}
	f := files[0]
	file, err = f.Open()
	if err != nil {
		return
	}
	return
}

func GetOrientation(f *os.File) (ot int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			ot = 1
		}
	}()
	x, err := exif.Decode(f)
	if err != nil {
		return
	}
	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return
	}
	ot = tag.Int(0)
	return
}

func ResizerHandler(w http.ResponseWriter, r *http.Request) {
	method, width, height, file, err := ParseResizerParams(r)
	if err != nil {
		glog.Warning(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	src, err := TempFile(cfg.TempDir, "src_", ".jpg")
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	srcName := src.Name()
	defer os.Remove(srcName)
	defer src.Close()

	if _, err := io.Copy(src, file); err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	src.Seek(0, 0)

	ot, err := GetOrientation(src)
	if err != nil {
		fmt.Println(err.Error())
		glog.Warning(err)
	}
	fmt.Println(ot)

	dst, err := TempFile(cfg.TempDir, "dst_", ".jpg")
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dstName := dst.Name()
	defer os.Remove(dstName)
	defer dst.Close()

	Resize(method, srcName, dstName, width, height, ot)

	http.ServeFile(w, r, dstName)
}
