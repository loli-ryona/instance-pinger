package main

import (
	js "encoding/json"
	fwk "instance-pinger/framework"

	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-ping/ping"
	mail "github.com/xhit/go-simple-mail/v2"
)

//Vars
var (
	inst fwk.Instances
	smtp fwk.Mail
)

//Load instances.json
func init() {
	instances, err := os.Open("instances.json")
	if err != nil {
		fmt.Println("Error loading instances file. Error: ", err)
		os.Exit(1)
	}

	if err = js.NewDecoder(instances).Decode(&inst); err != nil {
		fmt.Println("Error decoding instances. Error: ", err)
		os.Exit(1)
	}
}

//Main script
func main() {
	//Create Array for downed instances
	failedPings := 0
	instLength := []fwk.Instances{}
	downedInst := []string{
		"<h1>Instance Alert</h1>",
		"<p>The following instances have not<br>returned a ping in the last 5 minutes<br></p>",
	}

	//Wait Group shit
	var wg sync.WaitGroup
	wg.Add(len(instLength))

	//Ping IPs
	//Go Routine (im shit at these)
	fmt.Println("Starting pings")
	for i := 0; i < len(instLength); i++ {
		go func(i int) {
			//Pinging Instance
			fmt.Println("Pinging " + inst.Name[i])
			online, err := pingInstance(inst.Addr[i])
			if err != nil {
				fmt.Println(err.Error())
			}
			if online {
				fmt.Println("Instance " + inst.Name[i] + " returned status of: UP")
			} else {
				failedPings++
				downedInst = append(downedInst, inst.Name[i]+" - "+inst.Addr[i])
				fmt.Printf("instance" + inst.Name[i] + " returned status of: DOWN")
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("Finished pings")

	if failedPings >= 1 {
		fmt.Printf("Found %v downed instances, preparing alert mail.", failedPings)
		sendAlert()
	}
}

//Ping instance status function
//Returns true if packets are returned
//Returns false if packets arent returned
func pingInstance(addr string) (bool, error) {
	//Create pinger for address
	pinger, err := ping.NewPinger(addr)
	pinger.SetPrivileged(true)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	//Send and receive 4 bytes
	pinger.Count = 4
	err = pinger.Run()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	status := pinger.Statistics().PacketsRecv //Get amount of packets received

	if status == 0 {
		return false, err
	} else {
		return true, err
	}

}

//Send alert email if there is downed instances
//This uses xhit's go-simple-mail library
func sendAlert() {
	//Vars
	encryptionType := mail.EncryptionNone

	//Read mailBody.html
	body, err := os.ReadFile("mailBody.txt")
	if err != nil {
		fmt.Println("Error reading mailBody.txt. Error: ", err)
		os.Exit(1)
	}

	bodyStr := string(body)

	//Start mail client
	server := mail.NewSMTPClient()

	//Load SMTP settings
	mailSettings, err := os.Open("smtp.json")
	if err != nil {
		fmt.Println("Error loading smtp file. Error: ", err)
		os.Exit(1)
	}

	if err = js.NewDecoder(mailSettings).Decode(&smtp); err != nil {
		fmt.Println("Error decoding smtp settings. Error: ", err)
		os.Exit(1)
	}

	switch {
	case smtp.Encryption == "None":
		encryptionType = mail.EncryptionNone
	case smtp.Encryption == "SSL":
		encryptionType = mail.EncryptionSSL
	case smtp.Encryption == "SSLTLS":
		encryptionType = mail.EncryptionSSLTLS
	case smtp.Encryption == "STARTTLS":
		encryptionType = mail.EncryptionSTARTTLS
	case smtp.Encryption == "TLS":
		encryptionType = mail.EncryptionTLS
	}

	//Assign SMTP settings
	server.Host = smtp.Host
	server.Port = smtp.Port
	server.Username = smtp.Username
	server.Password = smtp.Password
	server.Encryption = encryptionType
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	//Connect to client
	smtpClient, err := server.Connect()

	if err != nil {
		fmt.Println("Expected nil, instead returned: ", err)
	}

	//Creating actual message
	email := mail.NewMSG()
	email.SetFrom(smtp.From).AddTo(smtp.To).SetSubject(smtp.Subject)
	email.SetBody(mail.TextHTML, bodyStr)
	email.AddAlternative(mail.TextPlain, smtp.Alternative)
	email.SetDate(time.Now().String())
	email.SetPriority(mail.PriorityHigh)

	if email.Error != nil {
		fmt.Println("Expected nil, instead returned: ", email.Error)
	}

	//Send email
	err = email.Send(smtpClient)
	email.GetError()
	if err != nil {
		fmt.Println("Expected nil, instead returned: ", err)
	}
}
