package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var client = &http.Client{}
var authenticated bool = false

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	//fmt.Printf("The URL is : %s\n", u.String())
	//fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	//fmt.Printf("The URL is : %s\n", u.String())
	//fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}

func setupClient() {
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	client.Jar = jar
}

//Given an authenticated client writes status.html to given path
func getStatus() {

	if !authenticated {
		auth(client)
	}

	say("Fetching progress page...", 2)
	resp, err := client.Get(penet + "/progress")
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	say("Progress page downloaded.", 1)

	say("Writing page to "+settings["statusPath"], 3)
	ioutil.WriteFile(settings["statusPath"], b, permissions)

	say(string(b), 5)

	//fmt.Println(string(b))
}

func auth(client *http.Client) {

	say("Authenticating...", 2)

	form := make(url.Values)
	form.Set("username", settings["username"])
	form.Set("password", settings["password"])
	form.Set("remember", "1")
	form.Set("login", "Login")

	// Authenticate
	_, err := client.PostForm(penet+"/login", form)
	if err != nil {
		fmt.Printf("Error Authenticating: %s", err)
	}

	say("Authenticated", 1)
	authenticated = true

}

//takes problem number and solution: submits answer online
func submit(problem int, solution string) (worked bool, message string) {
	if !authenticated {
		auth(client)
	}

	pname := strconv.Itoa(problem)
	theURL := penet + "/problem=" + pname

	say("Fetching Problem... "+pname, 2)
	resp, err := client.Get(theURL)
	say("Problem Downloaded.", 1)

	if err != nil {
		return false, "Fetching problem failed"
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	page := string(b)

	capStart := strings.Index(page, "<img src=\"captcha")
	if capStart == -1 {
		if strings.Contains(page, "Completed on ") || strings.Contains(page, "Go to the thread for") {
			say("Problem Already Completed", 1)
			answerStart := strings.Index(page, "Answer:")
			trunc := page[answerStart:]
			aStart := strings.Index(trunc, "<b>")
			aEnd := strings.Index(trunc, "</b>")
			correctAnswer := trunc[aStart+3 : aEnd]

			if correctAnswer == solution {
				return true, ""
			} else {
				return false, "Problem Solved: " + correctAnswer
			}

		} else {

			return false, "Can't find captcha in problem page."
		}
	}
	capEnd := strings.Index(page[capStart+10:], "\"")
	capURL := page[capStart+10 : capStart+10+capEnd]

	say("Downloading Captcha...", 2)
	resp, err = client.Get(penet + "/" + capURL)
	say("Captcha Downloaded.", 1)

	if err != nil {
		return false, "Fetching captcha failed."
	}

	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println("Cracking Captcha...")
	captcha := crackCap(b)

	fmt.Println("Captcha Solved:", captcha)

	form := make(url.Values)
	form.Set("guess_"+pname, solution)
	form.Set("confirm", captcha)

	fmt.Println("Submitting...")
	//Submit
	resp, err = client.PostForm(theURL, form)
	if err != nil {
		return false, "Trouble submitting solution"
	}

	b, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	page = string(b)

	if strings.Contains(page, "answer_wrong.png") {
		return false, "Wrong answer!"
	}

	if strings.Contains(page, "answer_correct.png") {
		return true, ""
	}

	if strings.Contains(page, "The confirmation code you entered was not valid") {

		return false, "Captcha Failed!"
	}

	return false, "wtf?"

}
