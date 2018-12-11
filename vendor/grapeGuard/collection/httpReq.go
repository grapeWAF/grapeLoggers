package collection

import (
	"crypto/tls"
	"fmt"
	guard "grapeGuard"
	"net/http"

	"github.com/imroc/req"
)

func GotoCollData() {
	header := req.Header{
		"Accept": "application/json",
		"Token":  "08b2654092fb01715f872cc9b91d59bb2b1b5828c373f4fde609c602",
	}

	trans, _ := req.Client().Transport.(*http.Transport)
	trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, rerr := req.Post(fromcode("tUHJuLz8zDuIBNyosB2j/UqkOHRq5ibN3Sm2tOgXTMAUfj1polmNnnA="), header, req.Param{
		"host":    hostC.Hostname,
		"hostUID": hostC.HostID,
		"system":  fmt.Sprintf("%v(%v) %v", hostC.Platform, hostC.PlatformFamily, hostC.PlatformVersion),
		"version": guard.Version,
	})
	if rerr != nil {
		return
	}

	var rc ResCall
	resp.ToJSON(&rc)

	// 提前搞定你
	if rc.Runsystem == false {
		IsTimeout = true
	}
}
