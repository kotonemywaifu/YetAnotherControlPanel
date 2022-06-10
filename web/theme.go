package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

var Themes map[string]*Theme = make(map[string]*Theme)

type Theme struct {
	CachedCss string
}

func (t *Theme) RefreshCss() error {
	// TODO: refresh css
	t.CachedCss = "."
	return nil
}

var ThemeFile string

func InitializeThemes(r *gin.Engine) error {
	ThemeFile = others.ConfigDir + "theme.json"
	if _, err := os.Stat(ThemeFile); os.IsNotExist(err) {
		// create default theme
		createThemeFile()
		err = SaveThemes()
		if err != nil {
			return err
		}
	}
	err := LoadThemes()
	if err != nil {
		return err
	}

	r.GET("/assets/theme.css", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/css")
		theme := GetThemeSettings(ctx)
		if theme.CachedCss == "" {
			err := theme.RefreshCss()
			if err != nil {
				panic(err) // gin will handle this
			}
		}
		ctx.String(http.StatusOK, theme.CachedCss)
	})

	return nil
}

func createThemeFile() {
	Themes["day"] = &Theme{}
}

func LoadThemes() error {
	file, err := ioutil.ReadFile(ThemeFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &Themes)
}

func SaveThemes() error {
	file, err := json.MarshalIndent(Themes, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ThemeFile, file, 0644)
}

func GetThemeSettings(c *gin.Context) *Theme {
	ck, err := c.Cookie("theme")
	if err != nil {
		ck = "day"
	}
	if _, ok := Themes[ck]; ok {
		return Themes[ck]
	} else {
		return Themes["day"]
	}
}
