package i18n

func loadChinese() *Locale {
	locale := loadEnglish() // use English as the base
	locale.Lang = "中文"

	// api.login
	locale.Api.Login.FailedTooManyTimes = "登入失敗，請稍後再試"
	locale.Api.Login.InvalidAccountHash = "無效的帳號哈希"
	locale.Api.Login.InvalidAccountCredentials = "無效的帳號或密碼"

	// page.login
	locale.Page.Login.UsernameField = "帳號"
	locale.Page.Login.PasswordField = "密碼"
	locale.Page.Login.LoginButton = "登入"

	return locale
}
