package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	command := "/home/taozhengkai/.pyenv/shims/pip"
	out, err := exec.Command(command, "freeze").Output()
	if err != nil {
		log.Fatal(err)
	}
	packages := strings.Split(string(out), "\n")
	// ignore last newline
	if packages[len(packages)-1] == "" {
		packages = packages[:len(packages)-1]
	}
	if len(packages) < 1 {
		fmt.Println("no packages listed")
		return
	}
	var args []string
	args = append(args, "uninstall", "-y")
	for _, p := range packages {
		args = append(args, strings.Split(p, "=")[0])
	}
	o, err := exec.Command(command, args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", o)
}
