package provider

var providerSmsCodes = []string{
	"Topolo",
	"Rond",
	"Kildy",
}

var providerMmsCodes = []string{
	"Topolo",
	"Rond",
	"Kildy",
}

var providerVoiceCallsCodes = []string{
	"TransparentCalls",
	"E-Voice",
	"JustPhone",
}

var providerEmailCodes = []string{
	"Gmail",
	"Yahoo",
	"Hotmail",
	"MSN",
	"Orange",
	"Comcast",
	"AOL",
	"Live",
	"RediffMail",
	"GMX",
	"Protonmail",
	"Yandex",
	"Mail.ru",
}

// IsRightCode проверяет является ли переданный код кодом из списка провайдеров
func IsRightCode(code, codeType string) bool {
	var provider []string
	switch codeType {
	case "Email":
		provider = providerEmailCodes
	case "VoiceCall":
		provider = providerVoiceCallsCodes
	case "Mms":
		provider = providerMmsCodes
	case "Sms":
		provider = providerSmsCodes
	}
	for _, v := range provider {
		if v == code {
			return true
		}
	}
	return false
}
