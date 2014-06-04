//Contains setting/file interaction code

package main

import (
	"io/ioutil"
	"os/user"
	"strconv"
	"strings"
)

var settings map[string]string

const (
	penet     = "http://projecteuler.net"
	probCount = 1000 //some number > #problems
)

func setSettings(path string, settingMap *map[string]string) {
	settings = make(map[string]string)

	settings["capPath"] = parsePath("~/.euler-tools/captcha/") //trailing slash!
	settings["knownPath"] = parsePath("~/.euler-tools/known.txt")
	settings["statusPath"] = parsePath("~/.euler-tools/status.html")
	settings["imageViewer"] = "eog"
	settings["debug"] = "3"

	settings["extensions"] = ".go.py.c"
	settings["exec.go"] = "go run %s/%s"
	settings["exec.py"] = "python2 %s/%s"
	settings["exec.c"] = "c99 %s/%s && ./a.out"

	//TODO: use parsePath to parse those settings which are paths
	say("Reading settings from file...", 1)
	fileSets := getData(path)
	for key, val := range fileSets {
		//settings from file overwrite defaults
		settings[key] = val
	}

}

//helper function for putData
func proccess(a string) []byte {
	out := make([]byte, 0)
	for i := 0; i < len(a); i++ {
		out = append(out, a[i])
	}
	return out
}

//Writes map data to file at given path. (colon separated)
func putData(path string, data map[string]string) {
	out := ""
	for i := 0; i < probCount; i++ {
		word := strconv.Itoa(i)
		if ans, ok := data[word]; ok {
			out += word + ":" + ans + "\n"
		}
	}

	ioutil.WriteFile(path, proccess(out), permissions)
}

//Reads colon separated data into map which is returned
func getData(path string) map[string]string {
	sets := eulerimport(path)

	out := make(map[string]string)

	for _, line := range sets {
		two := strings.SplitN(line, ":", 2)
		out[two[0]] = two[1]
	}

	return out
}

//Checks the given answer against the known-file.
func check(x int, ans string) (present, correct bool) {
	known := getData(settings["knownPath"])

	if rightAnswer, ok := known[strconv.Itoa(x)]; ok {
		if ans == rightAnswer {
			return true, true
		} else {
			return true, false
		}

	}

	return false, false

}

//Takes answer and adds it to file if not already present
func list(x int, ans string) {
	known := getData(settings["knownPath"])

	if _, ok := known[strconv.Itoa(x)]; ok {
		say("Answer already in list.", 1)
		return
	}

	say("Adding answer to list...", 2)
	known[strconv.Itoa(x)] = ans
	putData(settings["knownPath"], known)
	say("Answer added to list.", 1)
}

//Stolen from eulerlib. Should be merged eventually
func eulerimport(filename string) []string {
	// read whole the file
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var output []string

	currentline := ""

	for _, char := range b {
		if char == 10 {
			output = append(output, currentline)
			currentline = ""
		} else {
			currentline += string(char)
		}
	}

	if currentline != "" {
		output = append(output, currentline)
	}

	return output

}

func parsePath(path string) string {

	usr, _ := user.Current()
	dir := usr.HomeDir + "/"

	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir, 1)
	}
	return path

}
