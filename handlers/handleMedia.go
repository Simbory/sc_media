package handlers

import (
	"fmt"
	"github.com/Simbory/mego"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"
	"github.com/Simbory/sc_media/config"
)

func init() {
	mediaPrefixes := config.MediaPrefix()
	if len(mediaPrefixes) > 0 {
		for _, prefix := range mediaPrefixes {
			if len(prefix) == 0 {
				continue
			}
			mego.Get(strings.TrimRight(prefix, "/")+"/*pathInfo", handleMedia)
		}
	}
}

func handleMedia(ctx *mego.Context) interface{} {
	pathInfo := ctx.RouteParamString("pathInfo")
	var filePath = config.SrcDir() + pathInfo
	state, err := os.Stat(filePath)
	if err != nil || state.IsDir() {
		return nil
	}
	w := ctx.Request().URL.Query().Get("w")
	if len(w) == 0 {
		w = ctx.Request().URL.Query().Get("width")
	}
	h := ctx.Request().URL.Query().Get("h")
	if len(h) == 0 {
		h = ctx.Request().URL.Query().Get("height")
	}
	width, err := strconv.ParseUint(w, 0, 32)
	if err != nil {
		width = 0
	}
	height, err := strconv.ParseUint(h, 0, 32)
	if err != nil {
		height = 0
	}
	dotIndex := strings.LastIndex(pathInfo, ".")
	if dotIndex < 0 {
		return nil
	}
	ext := pathInfo[dotIndex:]
	pathWithoutExt := pathInfo[:dotIndex]
	if (width > 0 || height > 0) && resizable(ext) {
		filePath = resizeMedia(filePath, config.CacheDir() + pathWithoutExt, ext, uint(width), uint(height))
	}
	return mego.File(filePath, "")
}

func resizable(ext string) bool {
	lower := strings.ToLower(ext)
	return lower == ".jpg" || lower == ".jpeg" || lower == ".png" || lower == ".bmp"
}

func resizeMedia(filePath string, pathWithoutExt, ext string, width, height uint) string {
	newImgPath := fmt.Sprintf("%s/%dX%d%s", pathWithoutExt, width, height, ext)
	state, err := os.Stat(newImgPath)
	if err == nil {
		if state.IsDir() {
			return filePath
		} else {
			return newImgPath
		}
	}
	// get srcImage
	srcFile, err := os.Open(filePath)
	if err != nil {
		return filePath
	}
	srcImg := getImage(srcFile, ext)
	if srcImg == nil {
		srcFile.Close()
		return filePath
	}
	srcFile.Close()

	newDir := getDir(newImgPath)
	state, err = os.Stat(newDir)
	if err != nil {
		err = os.MkdirAll(newDir, 0777)
		if err != nil {
			return filePath
		}
	} else {
		if !state.IsDir() {
			return filePath
		}
	}
	newFile, err := os.Create(newImgPath)
	if err != nil {
		return filePath
	}
	defer newFile.Close()
	img := resize.Resize(width, height, srcImg, resize.Lanczos3)
	saveImage(newFile, img, ext)
	return newImgPath
}

func getImage(file *os.File, ext string) image.Image {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil
		}
		return img
	case ".png":
		img, err := png.Decode(file)
		if err != nil {
			return nil
		}
		return img
	case ".bmp":
		img, err := bmp.Decode(file)
		if err != nil {
			return nil
		}
		return img
	}
	return nil
}

func saveImage(file *os.File, newImg image.Image, ext string) {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		jpeg.Encode(file, newImg, nil)
	case ".png":
		png.Encode(file, newImg)
	case ".bmp":
		bmp.Encode(file, newImg)
	}
}
