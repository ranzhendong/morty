package dpimageupdate

import (
	"alert"
	"datastructure"
	"github.com/spf13/viper"
	"k8sapi"
	"net/http"
	"time"
)

func RollBack(r *http.Request, token *viper.Viper) (err error) {
	var (
		a               datastructure.Request
		f               [1]string
		bodyContentByte []byte
		content         string
	)
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiServer
	if err, bodyContentByte = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	// replace version and name
	if err, bodyContentByte = replaceResourceVersion(a, bodyContentByte); err != nil {
		return
	}

	//patch struck
	if err = k8sapi.APIServerPatch(a, bodyContentByte, token, ""); err != nil {
		MyErrorChan <- MyError{err}
		errors <- 1
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a, time.Duration(1))
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}
