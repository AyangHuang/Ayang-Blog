package main

import (
	"net/http"
	"os/exec"
)

func autoDeploy(w http.ResponseWriter, req *http.Request) {
	// 无加密，别搞我！求！求！
	// 不过这是内部接口，外部接口在nginx，我不告诉你！！！
	command := "./autodeploy.sh"
	cmd := exec.Command("/bin/bash", command)
	if err := cmd.Run(); err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
	//加密错误。。。不知道原因
	//if req.Method == http.MethodPost {
	//	// 获取请求报文的内容长度
	//	length := req.ContentLength
	//	// 新建一个字节切片，长度与请求报文的内容长度相同
	//	body := make([]byte, length)
	//	// 读取请求主体，并将具体内容读入 body 中
	//	req.Body.Read(body)
	//	// 获取环境变量token
	//	token := os.Getenv("SECRET_TOKEN_GITHUB_AUTO_DEPLOY")
	//	// 获取请求头的加密
	//	head := req.Header.Get("X-Hub-Signature-256")
	//	// hmac是Hash-based Message Authentication Code的简写，就是指哈希消息认证码
	//	h := hmac.New(sha256.New, []byte(token))
	//	h.Write(body)
	//	signature := h.Sum(nil)
	//	// 转成十六进制
	//	//signatureStr := hex.EncodeToString(signature)
	//	signatureStr := "sha256" + string(signature)
	//	if hmac.Equal([]byte(signatureStr), []byte(head)) {
	//		// 执行shell脚本
	//		exec.Command("autodeploy.sh")
	//		w.WriteHeader(200)
	//	}
	//	w.WriteHeader(400)
	//}
	//w.WriteHeader(404)
}
func main() {
	http.HandleFunc("/autodeploy", autoDeploy)
	http.ListenAndServe(":1314", nil)
}
