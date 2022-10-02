package main

import "os/exec"


func main() {
	command := "./autodeploy.sh"
	cmd := exec.Command("/bin/bash", command)
	if err := cmd.Run(); err == nil {
		println("success")
	} else {
		println("false")
	}

}












