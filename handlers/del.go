package handlers

import (
	"github.com/Simbory/mego"
	"os"
	"strings"
	"github.com/Simbory/sc_media/config"
)

func init() {
	mego.Post("/del", del)
}

func del(ctx *mego.Context) interface{} {
	// check access token
	if !tokenValid(ctx) {
		return mego.PlainText("0")
	}
	err := ctx.Request().ParseForm()
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
	// remove the old file
	state, err := os.Stat(fileName)
	if err == nil && !state.IsDir() {
		err = os.Remove(fileName)
		if err != nil {
			return mego.PlainText("0")
		}
	}
	return mego.PlainText("1")
}