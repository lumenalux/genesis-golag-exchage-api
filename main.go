package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	_ "strconv"
	"strings"
)

type Conf struct {
	SMTPUsername string `yaml:"smtp_username"`
	SMTPPassword string `yaml:"smtp_password"`
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     string `yaml:"smtp_port"`
	Port         string `yaml:"api_port"`
}

func (config *Conf) getConfig() *Conf {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func rate(w http.ResponseWriter, r *http.Request) {
	response := map[string]int{"rate": getBitcoinExchange()}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func getBitcoinExchange() int {
	resp, err := http.Get("https://minfin.com.ua/currency/converter/btc-c-uah/?val1=1")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	parsedHTMLString := string(body)

	firstStageRegExp, _ := regexp.Compile("Я получу.+1 UAH = 0 btc</span></label>")
	firstStageSubstring := firstStageRegExp.FindString(parsedHTMLString)

	secondStageRegExp, _ := regexp.Compile("value=\".+?\"")
	secondStageSubstring := secondStageRegExp.FindString(firstStageSubstring)

	BTCExchangeString := strings.Replace(secondStageSubstring, "value=", "", 1)
	BTCExchangeString = strings.ReplaceAll(BTCExchangeString, "\"", "")
	BTCExchangeString = strings.ReplaceAll(BTCExchangeString, " ", "")

	BTCExchange, err := strconv.Atoi(BTCExchangeString)
	if err != nil {
		log.Fatalln(err)
	}

	return BTCExchange
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if isEmailAdded(email) {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Invalid email value")
	} else {
		addEmail(email)
		w.WriteHeader(200)
		fmt.Fprintf(w, "Email added")
	}
}

func getEmails() []string {
	file, err := os.Open("emails.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	var emails []string

	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return emails
}

func isEmailAdded(email string) bool {
	for _, addedEmail := range getEmails() {
		if email == addedEmail {
			return true
		}
	}
	return false
}

func addEmail(email string) {
	file, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.WriteString(email + "\n")

	if err != nil {
		log.Fatal(err)
	}
}

func sendEmails(w http.ResponseWriter, r *http.Request) {
	var config Conf
	config.getConfig()

	from := config.SMTPUsername
	password := config.SMTPPassword

	to := getEmails()

	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	message := []byte(fmt.Sprintf("Поточний курс за 1 BTC складає %v грн.", getBitcoinExchange()))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Exchange wasn't sent due to internal server error!")

		log.Println(err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "BTC to UAH exchange rate was sended")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func handleRequests() {
	var config Conf
	config.getConfig()

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/rate", rate).Methods("GET")
	myRouter.HandleFunc("/subscribe", subscribe).Methods("POST")
	myRouter.HandleFunc("/sendEmails", sendEmails).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+config.Port, myRouter))
}

func main() {
	handleRequests()
}
