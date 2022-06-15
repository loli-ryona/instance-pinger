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
	smtp fwk.Mail
	inst fwk.Instances
)

//Load instances
func init() {
	instances, err := os.Open("json/instances.json")
	if err != nil {
		fmt.Println("Error loading instances file. Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully opened instances.json")

	if err = js.NewDecoder(instances).Decode(&inst); err != nil {
		fmt.Println("Error decoding instances. Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully decoded instances")
	fmt.Println("Initial load done")
}

//Main script
func main() {
	fmt.Println("Starting main script")

	//Create Array for downed instances
	failedPings := 0
	instLength := inst.Name
	fmt.Println(instLength)
	downedInst := []string{}
	fmt.Println(downedInst)

	//Wait Group shit
	var wg sync.WaitGroup
	wg.Add(len(instLength))

	//Ping IPs
	//Go Routine (im shit at these)
	fmt.Println("Starting pings")
	for i := 0; i < len(instLength); i++ {
		cum := i
		go func(cum int) {
			//Pinging Instance
			defer wg.Done()
			time.Sleep(time.Second)
			fmt.Println("Pinging " + inst.Name[cum])
			online, err := pingInstance(inst.Addr[cum])
			if err != nil {
				fmt.Println(err.Error())
			}
			if online {
				fmt.Println("Instance " + inst.Name[cum] + " returned status of: UP")
			} else {
				failedPings++
				downedInst = append(downedInst, "<p>"+inst.Name[cum]+" - "+inst.Addr[cum]+"</p><br>\n")
				fmt.Printf("Instance " + inst.Name[cum] + " returned status of: DOWN\n")
			}

		}(cum)
	}
	wg.Wait()

	fmt.Println(failedPings)
	fmt.Println("Finished pings")

	if failedPings >= 1 {
		fmt.Printf("Found %v downed instances, preparing alert mail.\n", failedPings)
		sendAlert(downedInst)
	} else {
		fmt.Println("Found 0 downed instances")
	}
}

//Ping instance status function
//Returns true if packets are returned
//Returns false if packets arent returned
func pingInstance(addr string) (bool, error) {
	//Create pinger for address
	pinger, err := ping.NewPinger(addr)
	
	//Set this to false if on linux and run the following command in cli:
	//sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
	pinger.SetPrivileged(true)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	
	pinger.Timeout = time.Second * 10

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
func sendAlert(downed []string) {
	//Vars
	encryptionType := mail.EncryptionNone

	//Build body structure
	bodyStr := `<h1>Instance Alert</h1><br>
<p>The following instances have not<br>returned a ping in the last 5 minutes</p><br>` + "\n"
	fmt.Println("This is the current body structure")
	fmt.Println(bodyStr)
	fmt.Println("Adding downed instances to structure")
	for i := 0; i < len(downed); i++ {
		bodyStr = bodyStr + downed[i]
	}
	fmt.Println("This the body structure after added downed instances")
	fmt.Println(bodyStr)

	//This is just for debugging html body errors
	//Ive noticed gmail duplicates the bodies downed instances
	/* os.Create("body.html")
	f, err := os.OpenFile("body.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error oppening body.html. Error: ", err)
	}
	if _, err := f.WriteString(bodyStr); err != nil {
		fmt.Println("Error writing to body.html. Error: ", err)
	}
	if err := f.Close(); err != nil {
		fmt.Println("Error closing body.html. Error: ", err)
	} */

	//Start mail client
	server := mail.NewSMTPClient()
	fmt.Println("Starting mail client")

	//Load SMTP settings
	mailSettings, err := os.Open("json/smtp.json")
	if err != nil {
		fmt.Println("Error loading smtp file. Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully loaded smtp.json")

	if err = js.NewDecoder(mailSettings).Decode(&smtp); err != nil {
		fmt.Println("Error decoding smtp settings. Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully loaded SMTP settings")

	switch {
	case smtp.Encryption == "None":
		encryptionType = mail.EncryptionNone
		fmt.Println("Detected Encryption Type: None")
	case smtp.Encryption == "SSL":
		encryptionType = mail.EncryptionSSL
		fmt.Println("Detected Encryption Type: SSL")
	case smtp.Encryption == "SSLTLS":
		encryptionType = mail.EncryptionSSLTLS
		fmt.Println("Detected Encryption Type: SSL/TLS")
	case smtp.Encryption == "STARTTLS":
		encryptionType = mail.EncryptionSTARTTLS
		fmt.Println("Detected Encryption Type: STARTTLS")
	case smtp.Encryption == "TLS":
		encryptionType = mail.EncryptionTLS
		fmt.Println("Detected Encryption Type: TLS")
	}

	//Assign SMTP settings
	fmt.Println("Pre-connection SMTP configuring")
	server.Host = smtp.Host
	server.Port = smtp.Port
	server.Username = smtp.Username
	server.Password = smtp.Password
	server.Encryption = encryptionType
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	fmt.Println("Pre-connection SMTP settings configured")

	//Connect to client
	fmt.Println("Connecting to client")
	smtpClient, err := server.Connect()

	if err != nil {
		fmt.Println("Expected nil, instead returned: ", err)
	}
	fmt.Println("Client Connected")

	//Creating actual message
	fmt.Println("Building Message")
	email := mail.NewMSG()
	email.SetFrom(smtp.From).AddTo(smtp.To).SetSubject(smtp.Subject)
	for i := 0; i < len(downed); i++ {
		bodyStr = bodyStr + downed[i]
	}
	email.SetBody(mail.TextHTML, bodyStr)
	//email.AddAlternative(mail.TextPlain, smtp.Alternative)
	//email.SetDate(time.Now().String())
	email.SetPriority(mail.PriorityHigh)

	if email.Error != nil {
		fmt.Println("Expected nil, instead returned: ", email.Error)
	}

	//Send email
	fmt.Println("Sending email")
	err = email.Send(smtpClient)
	email.GetError()
	if err != nil {
		fmt.Println("Expected nil, instead returned: ", err)
		fmt.Println("Email failed to send")
		os.Exit(1)
	}
	fmt.Println("Email sent successfully!")
	os.Exit(0)
}
