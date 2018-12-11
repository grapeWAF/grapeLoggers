package remaps

import (
	"crypto/tls"
	log "grapeLoggers/clientv1"
)

// 目前支持的CIPHERS列表
var SupportedCiphersMap = map[string]uint16{
	"ECDHE-ECDSA-AES256-GCM-SHA384":      tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-RSA-AES256-GCM-SHA384":        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-ECDSA-AES128-GCM-SHA256":      tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-RSA-AES128-GCM-SHA256":        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-ECDSA-WITH-CHACHA20-POLY1305": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"ECDHE-RSA-WITH-CHACHA20-POLY1305":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"ECDHE-RSA-AES256-CBC-SHA":           tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"ECDHE-RSA-AES128-CBC-SHA":           tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"ECDHE-ECDSA-AES256-CBC-SHA":         tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"ECDHE-ECDSA-AES128-CBC-SHA":         tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"RSA-AES256-CBC-SHA":                 tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"RSA-AES128-CBC-SHA":                 tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"ECDHE-RSA-3DES-EDE-CBC-SHA":         tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"RSA-3DES-EDE-CBC-SHA":               tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
}

// 默认支持的Ciphers列表
var defaultCiphers = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
}

// Map of supported curves
var supportedCurvesMap = map[string]tls.CurveID{
	"X25519": tls.X25519,
	"P256":   tls.CurveP256,
	"P384":   tls.CurveP384,
	"P521":   tls.CurveP521,
}

// List of all the curves we want to use by default.
//
// This list should only include curves which are fast by design (e.g. X25519)
// and those for which an optimized assembly implementation exists (e.g. P256).
// The latter ones can be found here: https://github.com/golang/go/tree/master/src/crypto/elliptic
var defaultCurves = []tls.CurveID{
	tls.X25519,
	tls.CurveP256,
}

const (
	minTLSVer = tls.VersionSSL30
	maxTLSVer = tls.VersionTLS12
)

func CreateFromXml(item *MapsXml) *tls.Config {
	if item.SSLConf == nil {
		return nil
	}

	log.Info("构建Tls Ver:",item.From , ",tls Info:",*item.SSLConf)
	newTlsConf := &tls.Config{}

	newTlsConf.MinVersion = item.SSLConf.MinVersion
	newTlsConf.MaxVersion = item.SSLConf.MaxVersion

	if len(item.SSLConf.CiphersVer) == 0 {
		return newTlsConf // 直接默认的加密吧
	}

	for _,keyn := range item.SSLConf.CiphersVer {
		if id,ok := SupportedCiphersMap[keyn]; ok {
			newTlsConf.CipherSuites = append(newTlsConf.CipherSuites,id)
		}
	}

	return newTlsConf
}

func GetTlsConfig(host string) *tls.Config {
	item := GetHttpsRemapItem(host)
	if item == nil {
		return nil
	}

	return item.TlsConfig
}