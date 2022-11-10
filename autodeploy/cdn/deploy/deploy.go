package deploy

import (
	"autodeploy/cdn"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	UPDATE = "update"
	ADD    = "add"
	MD     = ".md"
	HOST   = "https://ayang.ink"
)

var (
	CDN cdn.Cdn = &cdn.Client{}
)

func AutoDeploy(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	// 解密错误，why
	//if !verify(body, req) {
	//	w.WriteHeader(400)
	//	return
	//} else {
	//	w.WriteHeader(200)
	//}
	w.WriteHeader(200)
	command := "/home/ayang/site/Ayang-Blog/autodeploy/autodeploy.sh"
	cmd := exec.Command("/bin/bash", command)
	// 执行 shell 脚本文件
	if err := cmd.Run(); err == nil {
		// 缓存刷新
		cdnFunc(body)
	}
}

// false，why
func verify(body []byte, req *http.Request) bool {
	if req.Method == http.MethodPost {
		// 获取环境变量token
		//token := "KlBdopv3Xz0R7KRz6ogFL6rHXVKtpzzf"
		token := os.Getenv("SECRET_TOKEN_GITHUB_AUTO_DEPLOY")
		// 获取加密
		head := req.Header.Get("X-Hub-Signature-256")
		println(head)
		// hmac是Hash-based Message Authentication Code的简写，就是指哈希消息认证码
		h := hmac.New(sha256.New, []byte(token))
		h.Write(body)
		signature := "sha256=" + hex.EncodeToString(h.Sum(nil))
		if hmac.Equal([]byte(signature), []byte(head)) {
			return true
		}
	}
	return false
}

// cdn 的刷新、预热等操作
func cdnFunc(body []byte) {
	message := gjson.GetBytes(body, "head_commit.message").String()
	if strings.Contains(message, UPDATE) {
		update(body)
	} else if strings.Contains(message, ADD) {
		add(body)
	}
}

func update(body []byte) {
	modified := gjson.GetBytes(body, "head_commit.modified").Array()
	var strs []string
	for _, m := range modified {
		if str := judgeMD(m.String()); str != "" {
			str = HOST + str
			strs = append(strs, str)
		}
	}
	if len(strs) > 0 {
		CDN.Update(true, strs...)
	}
}

func add(body []byte) {
	modified := gjson.GetBytes(body, "head_commit.added").Array()
	var strs []string
	for _, m := range modified {
		if str := judgeMD(m.String()); str != "" {
			str = HOST + str
			strs = append(strs, str)
		}
	}
	if len(strs) > 0 {
		CDN.Add(true, strs...)
	}
}

// judgeMD 判断后缀是否未 “.MD”，并返回 /目录/
func judgeMD(str string) string {
	if strings.HasSuffix(str, MD) {
		return cut(str) + "/"
	} else {
		return ""
	}
}

// 截取有效目录，去掉“.MD”后缀，把“ ”替换成“-”，并把目录中大写改成小写，返回 /目录
func cut(str string) string {
	start := strings.LastIndex(str, "/")
	str = str[start : len(str)-3]
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.ToLower(str)
	return str
}
