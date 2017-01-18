package config

import (
	"io/ioutil"
	"encoding/xml"
	"strings"
	"mime"
)

type siteConfig struct {
	PathPrefix  string `xml:"path_prefix"`
	Token       string `xml:"token"`
	MediaPrefix string `xml:"media_prefix"`
	Mimes struct {
		List []struct{
			Ext         string    `xml:"ext,attr"`
			ContentType string `xml:"contentType,attr"`
		} `xml:"mime"`
	} `xml:"mimes"`
	CacheDir string	`xml:"cache_dir"`
	SrcDir   string `xml:"src_dir"`
}

var cfgData *siteConfig

func PathPrefix() string {
	if cfgData == nil {
		return ""
	}
	return cfgData.PathPrefix
}

func Token() string {
	if cfgData == nil {
		return ""
	}
	return cfgData.Token
}

func MediaPrefix() string {
	if cfgData == nil {
		return ""
	}
	return cfgData.MediaPrefix
}

func CacheDir() string {
	if cfgData == nil {
		return ""
	}
	return cfgData.CacheDir
}

func SrcDir() string {
	if cfgData == nil {
		return ""
	}
	return cfgData.SrcDir
}

func init() {
	var cfg = siteConfig{}
	cfgStr, err := ioutil.ReadFile("./config.xml")
	if err != nil {
		println("cannnot file the config file 'config.xml'")
		return
	}
	err = xml.Unmarshal([]byte(cfgStr), &cfg)
	if err != nil {
		println(err.Error())
		return
	}
	if len(cfg.PathPrefix) == 0 {
		println("Invalid config section: path_prefix cannot be empty.")
		return
	} else {
		cfg.PathPrefix = strings.Replace(cfg.PathPrefix, "\\", "/", -1)
	}
	if len(cfg.Token) == 0 {
		println("Invalid config section: token cannot be empty.")
		return
	}
	if len(cfg.MediaPrefix) == 0 {
		println("Invalid config section: media_prefix cannot be empty.")
		return
	}
	if len(cfg.CacheDir) == 0 {
		cfg.CacheDir = strings.Replace(cfg.PathPrefix+"/mediaCache/", "//", "/", -1)
	}
	if len(cfg.SrcDir) == 0 {
		cfg.SrcDir = strings.Replace(cfg.PathPrefix+"/mediaSrc/", "//", "/", -1)
	}
	if len(cfg.Mimes.List) > 0 {
		for _, mimeSetting := range cfg.Mimes.List {
			mime.AddExtensionType(mimeSetting.Ext, mimeSetting.ContentType)
		}
	}
	cfgData = &cfg
}