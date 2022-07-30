# **Docker BTC to UAH exchange API**
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

This is test task for Genesis Education. This is an API for viewing the bitcoin hryvnia exchange rate

## Installation
Open the directory where you want to put the project. Open a terminal on it
and put the command below:
```bash
git clone git@github.com:<YOUR_LOGIN>/genesis-golag-exchage-api.git exchangeApi
```

```bash
cd  exchangeApi
```
And then build docker image:
```bash
docker docker build --tag exchange-api .
```

## Usage
To run the application use the command below:
```bash
docker docker run -p 8080:8081 -p 587:587 -p 80:80 exchange-api
```

To get json with the bitcoin rate, you need to make a GET request
at the following address: ```<YOUR DOMAIN>:8080/rate```

To subscribe email, you need to make a POST request: 
```<YOUR DOMAIN>:8080/subscribe?email=subscriber@email.com```

To send emails to all subscribers, you need to make a POST request: 
```<YOUR DOMAIN>:8080/sendEmails```

## App settings
To use the application in full, you need to enter the settings of
your smtp server, from which messages for subscribed emails
addresses will be sent. Enter the required details in config.yaml
and then rebuild the docker image.

All settings are in config.yaml
```yaml
smtp_username: "test@domain.com"
smtp_password: "examplepassword"
smtp_host: "smtp.example.com"
smtp_port: "587"

api_port: "8081"
```

## PS
Yes.. I know that go the full name of the go language is "golang" not "golag". I just made a little mistake in the name of the repository when I did it in a hurry. Please forgive me)
