package i18n

func loadEnglish() *Locale {
	locale := &Locale{
		Lang: "English",
	}

	// api.login
	locale.Api.Login.FailedTooManyTimes = "failed to login too many times, please try again later"
	locale.Api.Login.InvalidAccountHash = "invalid account hash"
	locale.Api.Login.InvalidAccountCredentials = "invalid account credentials"

	return locale
}
