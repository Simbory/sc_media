package handlers

import (
	"github.com/Simbory/mego"
	"github.com/Simbory/sc_media/config"
)

func tokenValid(ctx *mego.Context) bool {
	accessToken := ctx.Request().Header.Get("Media-Upload-AccessToken")
	if len(config.Token()) == 0 || accessToken != config.Token() {
		return false
	}
	return true
}
