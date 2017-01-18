package handlers

import (
	"github.com/Simbory/mego"
	"io"
	"os"
	"strings"
	"github.com/Simbory/sc_media/config"
)

func init() {
	mego.Post("/upload", upload)
}

func upload(ctx *mego.Context) interface{} {
	// check access token
	if !tokenValid(ctx) {
		return mego.PlainText("0")
	}
	// Parse form
	err := ctx.Request().ParseMultipartForm(32 << 20)
	if err != nil {
		return mego.PlainText("0")
	}
	// get the form value
	var itemPath = ctx.Request().FormValue("path")
	var ext = ctx.Request().FormValue("ext")
	if len(itemPath) == 0 || len(ext) == 0 {
		return mego.PlainText("0")
	}
	fileName := strings.Replace(config.PathPrefix() + itemPath + "." + ext, "//", "/", -1)
	// get the uploaded file
	f, _, err := ctx.Request().FormFile("Filedata")
	if err != nil {
		return mego.PlainText("0")
	}
	defer f.Close()
	// remove the old file
	state, err := os.Stat(fileName)
	if err == nil && !state.IsDir() {
		err = os.Remove(fileName)
		if err != nil {
			return mego.PlainText("0")
		}
	}
	// check the dir, if the dir does not exist, create the dir
	dir := getDir(fileName)
	state, err = os.Stat(dir)
	if err != nil || !state.IsDir() {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return mego.PlainText("0")
		}
	}
	// create new empty file
	fw, err := os.Create(fileName)
	if err != nil {
		return mego.PlainText("0")
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		fw.Close()
		os.Remove(fileName)
		return mego.PlainText("0")
	} else {
		fw.Close()
	}
	return mego.PlainText("1")
}
