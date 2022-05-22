## Usage
**instances.json** is where any domains or ip addresses go that you wish to ping.<br>
**smtp.json** is where your mail providers smtp settings go, if you dont know them google "gmail smtp settings" or whatever your provider is.<br>

## SMTP Settings
Host - Mail provider host (e.g smtp.google.com)<br>
Port - Mail provider port (e.g 587, 465, 25, etc)<br>
Username - Auth Username<br>
Password - Auth Password<br>
Encryption - Type of security to use (None, SSL, SSLTLS, STARTTLS, TLS)<br>
From - From message<br>
To - Email to send to<br>
Subject - Subject of the email<br>
Alternative - Alt text if enabled in main.go (not needed usually)<br>

