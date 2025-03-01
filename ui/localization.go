package ui

import (
	"syscall"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/leonelquinteros/gotext"
	"github.com/lxn/win"
)

func GetLocalization() *gotext.Po {
	locale := getUserLocale()
	return loadPO(locale)
}

func loadPO(locale string) *gotext.Po {
	po := gotext.NewPoFS(assets.Locales)
	_, err := assets.Locales.ReadDir("locales/" + locale)
	if err != nil {
		locale = "en"
	}
	po.ParseFile("locales/" + locale + "/default.po")
	return po
}

func getUserLocale() string {
	var buf [4]uint16
	win.GetLocaleInfo(win.LOCALE_USER_DEFAULT, win.LOCALE_SISO639LANGNAME, &buf[0], int32(len(buf)))

	return syscall.UTF16ToString(buf[:])
}
