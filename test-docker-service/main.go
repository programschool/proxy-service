package main

import (
	"fmt"
	"github.com/programschool/proxy-service/library/dockertools"
	"strings"
)

func main() {
	exec()
}

func exec() {
	server := "192.168.50.104"
	dockerHost := fmt.Sprintf("%s:%s", server, "2376")

	dock := dockertools.Dock{}.New(dockerHost)
	defer dock.Close()

	containerID := "3f5859de8f84"
	code := "I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAvLyBwcmludGYoKSBkaXNwbGF5cyB0aGUgc3RyaW5nIGluc2lkZSBxdW90YXRpb24KICAgcHJpbnRmKCJIZWxsbywgV29ybGQhIik7CiAgIHJldHVybiAwOwp9Cg=="
	fileName := "aGVsbG8uYwo="
	execCommand := "Z2NjIC1vIGJpbi9oZWxsbwo="
	callBackFile := "YmluL2hlbGxvCg=="

	command := strings.Split(
		fmt.Sprintf("bash run.sh %s %s %s %s", execCommand, code, fileName, callBackFile),
		" ",
	)
	dock.ExecCommand("/programschool/execute", containerID, command)
	//fmt.Println("err")
	//fmt.Println(err)
	//fmt.Println("inspect")
	//fmt.Println(inspect)

	// gcc file params exec
}

//func main() {
//	file, err := "/data/test.txt", "file not found"
//
//	log("File {file} had error {error}", "{file}", file, "{error}", err)
//}
//
//func log(format string, args ...string) {
//	r := strings.NewReplacer(args...)
//	fmt.Println(r.Replace(format))
//}
