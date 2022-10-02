package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"os/exec"
	"fmt"
)

func autoDeploy(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		// 获取请求报文的内容长度
		length := req.ContentLength
		// 新建一个字节切片，长度与请求报文的内容长度相同
		body := make([]byte, length)
		// 读取请求主体，并将具体内容读入 body 中
		req.Body.Read(body)
		// 获取环境变量token
		//token := os.Getenv("SECRET_TOKEN_GITHUB_AUTO_DEPLOY")
		token := "KlBdopv3Xz0R7KRz6ogFL6rHXVKtpzzf"
		// 获取请求头的加密
		head := req.Header.Get("X-Hub-Signature-256")
		// hmac是Hash-based Message Authentication Code的简写，就是指哈希消息认证码
		h := hmac.New(sha256.New, []byte(token))
		h.Write(body)
		signature := h.Sum(nil)
		for _, v:= range signature {
			fmt.Printf("%c", v)
		}
		// 转成十六进制
		//signatureStr := hex.EncodeToString(signature)
		var pre []byte = []byte("sha256=")
		end := append(pre, signature...)
		for _, v := range end {
			fmt.Printf("%c", v)
		}
		if hmac.Equal(end, []byte(head)) {
			// 执行shell脚本
			exec.Command("autodeploy.sh")
			w.WriteHeader(200)
		}
		w.WriteHeader(400)
	}
	w.WriteHeader(400)
}
func main() {
	http.HandleFunc("/autodeploy", autoDeploy)
	http.ListenAndServe(":1314", nil)
}

