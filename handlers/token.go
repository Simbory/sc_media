package handlers

import (
	"github.com/Simbory/mego"
	"github.com/Simbory/sc_media/config"
)

func tokenValid(ctx *mego.Context) bool {
	if ctx == nil {
		return false
	}
	if len(config.Token()) == 0 {
		return true
	}
	accessToken := ctx.Request().Header.Get("Media-Upload-AccessToken")
	return config.Token() == accessToken
}
