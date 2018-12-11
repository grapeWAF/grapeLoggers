package jsGuard

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	jsGuardKey = "2fc0e78f129d6b306b7e02839d38a0329e476271"

	ttl = time.Second * 45 // 有效期
)

var (
	/*
			 if(browser.Chrome !== undefined) {
		                    barVars = "_c" + browser.Chrome;
		                }else if(browser.Edge !== undefined) {
		                    barVars = "_e" + browser.Edge;
		                }else if(browser.Gecko !== undefined) {
		                    barVars = "_g" + browser.Gecko;
		                }else if(browser.MSIE !== undefined) {
		                    barVars = "_j" + browser.MSIE + browser.rv;
		                }else if(browser.Opera !== undefined) {
		                    barVars = "_o" + browser.Opera;
		                }else if(browser.Safari !== undefined) {
		                    barVars = "_s" + browser.Safari;
		                }else if(browser.Webkit !== undefined) {
		                    barVars = "_w" + browser.Webkit;
		                }

		                if(browser.Windows !== undefined) {
		                    barVars += "_W" + browser.Windows;
		                }else if(browser.Android !== undefined) {
		                    barVars += "_A" + browser.Android;
		                }else if(browser.Macintosh !== undefined) {
		                    barVars += "_M" + browser.Macintosh;
		                }else if(browser.IOS !== undefined) {
		                    barVars += "_O" + browser.Android;
		                }
	*/
	badValue   = `evua41.556.222`
	browerType = map[uint8]string{
		'c': "Chrome ",
		'e': "Edge ",
		'g': "Firefox ",
		'j': "IE ",
		'o': "Opera ",
		's': "Safari ",
		'w': "Webkit ",
	}

	platformType = map[uint8]string{
		'W': "Windows ",
		'A': "Android ",
		'M': "Macintosh ",
		'O': "IOS ",
	}
)

func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func getBrowerType(src uint8) string {
	val, has := browerType[src]
	if !has {
		return "error"
	}

	return val
}

func getPlatformType(src uint8) string {
	val, has := platformType[src]
	if !has {
		return "error"
	}

	return val
}

func SplitBrower(key string) (rKey, brower, platform string) {
	brower = "unknow"
	platform = "unknow"
	pos := strings.Index(key, "_")
	if pos == -1 {
		rKey = key
		return
	}

	rKey = key[:pos]
	platData := key[pos+1:]
	if platData == badValue {
		return
	}
	pos = strings.Index(platData, "_")
	if pos == -1 {
		brower = getBrowerType(platData[0]) + platData[1:]
		return
	}

	brower = getBrowerType(platData[0]) + platData[1:pos]
	if pos+2 >= len(platData) {
		return
	}

	platform = getPlatformType(platData[pos+1]) + platData[pos+2:]
	return
}

func CreateHash(host, ip string, xorKey byte) (key string, en string) {
	timeUnix := time.Now().UnixNano() // 保证不会产生HASH对撞
	hash := hmac.New(sha1.New, []byte(jsGuardKey))
	ipInt := uint32(InetAtoN(ip) ^ 0x2f1256)
	hash.Write([]byte(fmt.Sprintf("%s_%x_%x", host, timeUnix, ipInt)))

	hashStr := hex.EncodeToString(hash.Sum(nil))
	keydatas := []byte(hashStr)
	for i := 0; i < len(keydatas); i++ {
		keydatas[i] = keydatas[i] ^ xorKey
	}

	key = hashStr
	en = base64.StdEncoding.EncodeToString(keydatas)
	return
}

func Packup(values url.Values) string {
	sb := strings.Builder{}
	for k, v := range values {
		prefix := k + "="
		for _, val := range v {
			if sb.Len() > 0 {
				sb.WriteByte('&')
			}
			sb.WriteString(prefix)
			sb.WriteString(val)
		}
	}
	return sb.String()
}

func RemoveKey(src, key string) string {
	pos := strings.LastIndex(src, key)
	if pos != -1 {
		if pos == 0 {
			return ""
		}
		if pos > 0 {
			pos--
		}
		return src[:pos]
	}

	return src
}
