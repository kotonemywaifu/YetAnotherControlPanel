package web

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

var Themes map[string]*Theme = make(map[string]*Theme)

type Theme struct {
	Primary        string `json:"primary"`
	PrimaryFocus   string `json:"primary-focus"`
	PrimaryContent string `json:"primary-content"`

	Secondary        string `json:"secondary"`
	SecondaryFocus   string `json:"secondary-focus"`
	SecondaryContent string `json:"secondary-content"`

	Accent        string `json:"accent"`
	AccentFocus   string `json:"accent-focus"`
	AccentContent string `json:"accent-content"`

	Neutral        string `json:"neutral"`
	NeutralFocus   string `json:"neutral-focus"`
	NeutralContent string `json:"neutral-content"`

	Base100     string `json:"base-100"`
	Base200     string `json:"base-200"`
	Base300     string `json:"base-300"`
	BaseContent string `json:"base-content"`

	Info    string `json:"info"`
	Success string `json:"success"`
	Warning string `json:"warning"`
	Error   string `json:"error"`

	RoundedBox      string `json:"rounded-box"`
	RoundedButton   string `json:"rounded-button"`
	RoundedBadge    string `json:"rounded-badge"`
	AnimationButton string `json:"animation-button"`
	AnimationInput  string `json:"animation-input"`
	PaddingCard     string `json:"padding-card"`
	ButtonTextCase  string `json:"button-text-case"`
	NavbarPadding   string `json:"navbar-padding"`
	BorderButton    string `json:"border-button"`
	FocusRing       string `json:"focus-ring"`
	FocusRingOffset string `json:"focus-ring-offset"`

	CachedCss string `json:"-"`
}

func (t *Theme) GetTemplated() (string, error) {
	tmpl := template.Must(template.New("Theme").Parse(util.Must(util.ReadFile(templatesFS, "templates/theme.css"))))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, t)
	if err != nil {
		return "", err
	}
	res := buf.String()
	if others.TheConfig.MinifyResources {
		res, err = util.MinifyCss(res)
	}
	return res, err
}

func (t *Theme) RefreshCss() error {
	css, err := t.GetTemplated()
	if err != nil {
		return err
	}
	t.CachedCss = css
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
		var css string
		if others.TheConfig.CacheTemplate {
			if theme.CachedCss == "" {
				err := theme.RefreshCss()
				if err != nil {
					panic(err) // gin will handle this
				}
			}
			css = theme.CachedCss
		} else {
			css, err = theme.GetTemplated()
			if err != nil {
				panic(err) // gin will handle this
			}
		}
		ctx.String(http.StatusOK, css)
	})

	return nil
}

func createThemeFile() {
	// themes from daisyUI https://github.com/saadeghi/daisyui/blob/d84501e5aa67bac6c432e5a01eefcc685937a789/src/themes/
	basicTheme := Theme{
		Primary:          "259 94% 51%",
		PrimaryFocus:     "259 94% 41%",
		PrimaryContent:   "0 0% 100%",
		Secondary:        "314 100% 47%",
		SecondaryFocus:   "314 100% 37%",
		SecondaryContent: "0 0% 100%",
		Accent:           "174 60% 51%",
		AccentFocus:      "174 60% 41%",
		AccentContent:    "0 0% 100%",
		Neutral:          "219 14% 28%",
		NeutralFocus:     "222 13% 19%",
		NeutralContent:   "0 0% 100%",
		Base100:          "0 0% 100%",
		Base200:          "210 20% 98%",
		Base300:          "216 12% 84%",
		BaseContent:      "215 28% 17%",
		Info:             "207 90% 54%",
		Success:          "174 100% 29%",
		Warning:          "36 100% 50%",
		Error:            "14 100% 57%",
		RoundedBox:       "1rem",
		RoundedButton:    "0.5rem",
		RoundedBadge:     "9999px",
		AnimationButton:  "0.25s",
		AnimationInput:   "0.4s",
		PaddingCard:      "2rem",
		ButtonTextCase:   "uppercase",
		NavbarPadding:    "0.5rem",
		BorderButton:     "1px",
		FocusRing:        "2px",
		FocusRingOffset:  "2px",
	}
	// based on basic theme
	dayTheme := basicTheme // copy struct
	Themes["day"] = &dayTheme

	nightTheme := basicTheme
	nightTheme.Primary = "259 94% 61%"
	nightTheme.PrimaryFocus = "259 94% 51%"
	nightTheme.Neutral = "223 13% 19%"
	nightTheme.NeutralFocus = "223 14% 10%"
	nightTheme.Base100 = "219 14% 28%"
	nightTheme.Base200 = "222 13% 19%"
	nightTheme.Base300 = "223 14% 10%"
	nightTheme.BaseContent = "228 14% 93%"
	Themes["night"] = &nightTheme
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
