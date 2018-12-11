////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/10/19
//  执行模版类
////////////////////////////////////////////////////////////
package grapeGuard

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"path"
	"strings"

	log "grapeLoggers/clientv1"
)

var (
	// 構建一組頁面，請求錯誤時返回。
	templates map[string]*template.Template = map[string]*template.Template{}
)

func loadTemplates() {
	files, err := ioutil.ReadDir("templates")
	if err != nil {
		log.Error("载入模版错误:", err)
		return
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		ext := path.Ext(v.Name())
		baseName := strings.TrimSuffix(v.Name(), ext)
		if strings.ToLower(ext) != ".html" {
			continue
		}

		log.Info("加载模版文件:", v.Name())
		tp, terr := template.ParseFiles("templates/" + v.Name())
		if terr != nil {
			log.Error("载入模版错误:", terr, ",File:", v.Name())
			continue
		}
		templates[baseName] = tp
	}
}

func ExecutTP(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	tp, ok := templates[fmt.Sprint(status)]
	if !ok {
		return
	}

	tp.Execute(w, nil)
}
