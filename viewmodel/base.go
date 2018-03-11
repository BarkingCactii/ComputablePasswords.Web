package viewmodel

type Base struct {
	Title     string
	WebURL    string
	Result    string
	PublicKey string
}

func NewBase() Base {
	return Base{
		Title:     "Computable Passwords",
		WebURL:    "Enter URL",
		Result:    "",
		PublicKey: "nameddmmyy for exanple",
	}
}
