package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ExampleResponse struct {
	Answer string `json:"answer"`
}

type Rules struct {
	Number int    `json:"number"`
	Resp   string `json:"response"`
}

type Response struct {
	Message string          `json:"message"`
	Next    string          `json:"nextQuestion"`
	Numbers []int           `json:"numbers"`
	Result  string          `json:"result"`
	Rules   []Rules         `json:"rules"`
	Example ExampleResponse `json:"exampleResponse"`
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func checkAnswer(url, ans string) (Response, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(ExampleResponse{Answer: ans})
	checkError(err)
	res, err := http.Post(url, "application/json", buf)
	checkError(err)

	return parseResponse(res)
}

func separator() {
	fmt.Println("\n----------------------------------------------------------\n")
}

func displayQuestion(resp Response) {
	fmt.Println(resp.Message)

	if resp.Rules != nil {
		fmt.Println("\nRules:")
		for _, res := range resp.Rules {
			fmt.Printf("    For %d print %s\n", res.Number, res.Resp)
		}
	}

	fmt.Println("\nExample response:", resp.Example.Answer)

	if resp.Numbers != nil {
		fmt.Println("\nNumbers:", resp.Numbers)
	}

	separator()
}

func getQuestion(url string) (Response, error) {
	res, err := http.Get(url)
	checkError(err)

	return parseResponse(res)
}

func parseResponse(res *http.Response) (Response, error) {
	var resp Response
	body, err := ioutil.ReadAll(res.Body)

	checkError(err)
	err = json.Unmarshal(body, &resp)
	checkError(err)

	return resp, nil
}

func solve(nums []int, rules []Rules) string {
	st := ""

	for _, num := range nums {
		temp := ""
		for _, rule := range rules {
			if num%rule.Number == 0 {
				temp += rule.Resp
			}
		}
		if temp == "" {
			temp = strconv.Itoa(num)
		}
		st += temp + " "
	}

	return st
}

func interactive(domain string, resp Response) {
	var ans string
	next := resp.Next

	fmt.Println("-------------------- Starting fizzbot --------------------\n")
	fmt.Println("Dear Candidate\n")
	fmt.Print(resp.Message)
	separator()

	resp, err := getQuestion(domain + next)
	checkError(err)
	displayQuestion(resp)

	for {
		if resp.Numbers == nil {
			ans = "go"
		} else {
			ans = solve(resp.Numbers, resp.Rules)
		}
		fmt.Println(ans)
		resp, err = checkAnswer(domain+next, strings.Trim(ans, "\r\n "))
		checkError(err)
		if strings.Trim(resp.Result, "\r\n ") == "correct" {
			separator()
			fmt.Println(resp.Message)
			separator()
			next = resp.Next
			resp, err = getQuestion(domain + next)
			checkError(err)
			displayQuestion(resp)
		} else if strings.Trim(resp.Result, "\r\n ") == "interview complete" {
			separator()
			fmt.Print(resp.Message)
			separator()
			break
		} else {
			fmt.Println(resp.Message)
		}
	}
}

func main() {
	var resp Response

	domain := "https://api.noopschallenge.com"
	res, err := http.Get(domain + "/fizzbot")
	checkError(err)
	resp, err = parseResponse(res)
	checkError(err)

	interactive(domain, resp)
}
