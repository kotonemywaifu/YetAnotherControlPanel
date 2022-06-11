package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Locale struct {
	Lang string
	Api  struct {
		Login struct {
			FailedTooManyTimes        string
			InvalidAccountHash        string
			InvalidAccountCredentials string
		}
	}
	Page struct {
		Login struct {
			UsernameField string
			PasswordField string
			LoginButton   string
		}
	}
}

var locales map[string]*Locale = make(map[string]*Locale)

func GetLocale(lang string) *Locale {
	lang = strings.ToLower(lang) // don't care about case

	if val, ok := locales[lang]; ok {
		return val
	} else {
		// load locale dynamically to save memory
		locale := loadLocale(lang)
		locales[lang] = locale
		return locale
	}
}

func loadLocale(lang string) *Locale {
	switch lang {
	case "en":
		return loadEnglish()
	case "zh":
		return loadChinese()
	default:
		return GetLocale("en") // fallback to that pointer
	}
}

func ReadLocale(c *gin.Context) *Locale {
	if c.GetHeader("Accept-Language") != "" {
		lang := strings.Split(c.GetHeader("Accept-Language"), ",")[0]
		lang = strings.Split(lang, ";")[0]
		lang = strings.Split(lang, "-")[0]
		return GetLocale(lang)
	} else if c, err := c.Request.Cookie("language"); err == nil {
		return GetLocale(c.Value)
	} else {
		return GetLocale("en")
	}
}
