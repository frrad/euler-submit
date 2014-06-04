package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func detectExec(n int, files []string) (filename, extension string) {
	extensions := strings.Split(settings["extensions"], ".")

	//take care of leading .
	for len(extensions[0]) == 0 {
		extensions = extensions[1:]
	}

	//fmt.Println(extensions, len(extensions))

	for _, extension := range extensions {

		expr := "Problem" + "0*" + strconv.Itoa(n) + "\\." + extension

		re, _ := regexp.Compile(expr)

		for _, file := range files {
			if re.MatchString(file) {
				return file, extension
			}
		}
	}

	return "", ""
}

func resolveCmd(path, name, extension string) []string {
	pfstr := settings["exec."+extension]
	//fmt.Println(pfstr)

	mid := fmt.Sprintf(pfstr, path, name)

	//fmt.Println(mid)

	return strings.Split(mid, " ")
}

func runProb(n int) (works bool, message, output string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	files, err := dir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	nstr, ext := detectExec(n, files)
	if ext == "" {
		return false, "Can't find file.", ""
	}

	resolution := resolveCmd(path, nstr, ext)
	//fmt.Println(resolution)
	cmd := exec.Command(resolution[0], resolution[1:]...)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return false, "Trouble getting pipe.", ""
	}
	if err := cmd.Start(); err != nil {
		return false, "Trouble starting program.", ""
	}

	b, err := ioutil.ReadAll(out)
	if err != nil {
		return false, "Trouble reading output.", ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		return false, "Program exitted with error.", ""
	}

	out.Close()

	programOutput := strings.Split(string(b), "\n")

	time := ""

	for i := len(programOutput) - 1; i >= 0; i-- {

		line := programOutput[i]

		say(line, 5)

		if strings.Contains(line, "Elapsed") {
			time = line

		}
		if len(line) > 0 && !strings.Contains(line, "Elapsed") {
			return true, time, line

		}
	}

	return false, "Can't find output", ""

}
