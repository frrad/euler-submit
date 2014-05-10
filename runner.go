package main

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

func runProb(n int) (works bool, message, output string) {
	nstr := strconv.Itoa(n)
	for len(nstr) < 3 {
		nstr = "0" + nstr
	}
	nstr = "Problem" + nstr

	cmd := exec.Command("go", "run", nstr+".go")
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
