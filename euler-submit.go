package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const permissions = 0666

func say(message string, level int) {
	debugLevel, _ := strconv.Atoi(settings["debug"])
	if debugLevel >= level {
		fmt.Println(message)
	}
}

func crackCap(b []byte) (crack string) {
	timeStr := strconv.FormatInt((time.Now().Unix()), 10)
	path := settings["capPath"] + timeStr + ".png"

	ioutil.WriteFile(path, b, permissions)

	do := exec.Command(settings["imageViewer"], path)
	do.Start()

	fmt.Println("Please Input Captcha:")
	fmt.Scan(&crack)

	return
}

//A wrapper for submit which checks against known list
func fancySubmit(x int, ans string) bool {

	if known, correct := check(x, ans); !known {

		if worked, mess := submit(x, ans); worked {
			say("Correct!", 0)
			say(penet+"/thread="+strconv.Itoa(x), 1)

			list(x, ans)

			getStatus()

			return true

		} else if len(mess) > 14 && mess[:14] == "Problem Solved" {
			say("Wrong answer! (Problem Already Solved)", 0)

			list(x, mess[16:])

		} else {
			say(mess, 0)
		}

	} else {
		say("Answer in list:", 1)
		if correct {
			say("Correct!", 0)
			say(penet+"/thread="+strconv.Itoa(x), 1)

		} else {
			say("Wrong answer!", 0)

		}
	}

	return false

}

//Parses user-submitted problem range specification, returning a list of
//problems, and err == nil if successful.
func parse(spec string) (list []int, err bool) {
	splitted := strings.Split(spec, "-")

	if len(splitted) > 2 || len(splitted) < 1 {
		return nil, true
	}

	if len(splitted) == 1 {
		if pnumber, err := strconv.Atoi(os.Args[1]); err == nil {
			return []int{pnumber}, false
		} else {
			return nil, true
		}
	}

	a, b := 1, 1

	if splitted[0] == "" {
		if temp, err := strconv.Atoi(splitted[1]); err == nil {
			b = temp
		} else {
			return nil, true
		}
	} else {
		tempa, err1 := strconv.Atoi(splitted[0])
		tempb, err2 := strconv.Atoi(splitted[1])
		if err1 != nil || err2 != nil {
			return nil, true
		}

		a, b = tempa, tempb

	}

	list = make([]int, b-a+1)
	for i := range list {
		list[i] = i + a
	}

	return
}

func main() {
	setupClient()

	setPath := parsePath("~/.euler-tools/settings.dat")
	setSettings(setPath, &settings)

	switch len(os.Args) {
	case 1:
		say("No arguments!", 0)

	case 2:
		//Only one argument
		switch argue := os.Args[1]; argue {

		case "R":
			say("Updating Status:", 0)
			getStatus()

		default:
			if plist, err := parse(argue); err == false {

				for _, pnumber := range plist {

					say("Solving #"+strconv.Itoa(pnumber), 1)

					if works, mess, out := runProb(pnumber); works {
						say("Answer: "+out, 1)

						if mess != "" { //time
							say(mess, 2)
						}

						fancySubmit(pnumber, out)

					} else {
						fmt.Println(mess)
					}

					fmt.Print("\n")
				}

			} else {
				say("Can't parse argument!", 0)
			}
		}
	case 3:
		if pnumber, err := strconv.Atoi(os.Args[1]); err == nil {
			out := os.Args[2]
			say("Submitting: "+out, 1)
			fancySubmit(pnumber, out)
		} else {
			say("Can't parse problem number!", 0)
		}
	default:
		say("Too many arguments!", 0)

	}

}
