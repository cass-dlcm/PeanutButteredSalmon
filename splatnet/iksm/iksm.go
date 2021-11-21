package iksm

import (
	"github.com/cass-dlcm/splatnetiksm"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// GenNewCookie attempts to generate a new cookie in case the provided one is invalid.
func GenNewCookie(reason string, _ string, client *http.Client) {
	sessionToken, cookie, errs := splatnetiksm.GenNewCookie(viper.GetString("user_lang"), viper.GetString("session_token"), reason, client)
	if len(errs) > 0 {
		log.Panicln(errs)
	}
	viper.Set("session_token", sessionToken)
	viper.Set("cookie", cookie)
	if err := viper.WriteConfig(); err != nil {
		log.Panicln(err)
	}
}
