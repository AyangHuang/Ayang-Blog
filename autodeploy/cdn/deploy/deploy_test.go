package deploy

import (
	"fmt"
	"testing"
)

const (
	bodyUpdate = `{"head_commit": {
		"id": "e81e9edb357166bd13adad56e4d592f45ab65fc6",
		"tree_id": "685f775c02da331b69ea35b6a0685ef6978a4f01",
		"distinct": true,
		"message": "update post",
		"timestamp": "2022-11-07T00:29:24+08:00",
		"url": "https://github.com/AyangHuang/Ayang-Blog/commit/e81e9edb357166bd13adad56e4d592f45ab65fc6",
		"author": {
		  "name": "ONE-YANG",
		  "email": "3392406201@qq.com",
		  "username": "AyangHuang"
		},
		"committer": {
		  "name": "ONE-YANG",
		  "email": "3392406201@qq.com",
		  "username": "AyangHuang"
		},
		"added": [
			
		],
		"removed": [
	
		],
		"modified": [
		  "content/posts/底层_通过 go 汇编了解函数调用栈帧.md",
		  "content/posts/底层_通"
		]
	  }}`
	bodyAdd = `{"head_commit": {
		"id": "e81e9edb357166bd13adad56e4d592f45ab65fc6",
		"tree_id": "685f775c02da331b69ea35b6a0685ef6978a4f01",
		"distinct": true,
		"message": "add post",
		"timestamp": "2022-11-07T00:29:24+08:00",
		"url": "https://github.com/AyangHuang/Ayang-Blog/commit/e81e9edb357166bd13adad56e4d592f45ab65fc6",
		"author": {
		  "name": "ONE-YANG",
		  "email": "3392406201@qq.com",
		  "username": "AyangHuang"
		},
		"committer": {
		  "name": "ONE-YANG",
		  "email": "3392406201@qq.com",
		  "username": "AyangHuang"
		},
		"added": [
			"content/posts/add底层_通过 go 汇编了解函数调用栈帧.md",
		  	"content/posts/add底层_通"
		],
		"removed": [
	
		],
		"modified": [

		]
	  }}`
)

func TestJudgeMD(t *testing.T) {
	str := "content/posts/底层_通过 go 汇编了解函数调用栈帧.MD"
	str = judgeMD(str)
	fmt.Print(str)
}

func TestAutoDeploy(t *testing.T) {
	update([]byte(bodyUpdate))
	add([]byte(bodyAdd))
}
