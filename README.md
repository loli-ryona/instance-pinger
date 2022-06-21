# Instance Pinger
This is a ip/domain pinger made in GoLang that sends an alert email on failed pings.<br>
I primarily use this for keeping track of my AWS Lightsail instances as they go down often.<br>

## Usage
In the **json** folder is where your settings are located.<br>
The **smtp.json** contains the smtp settings. Put your mail providers settigns in here.<br>
The **instances.json** is where your domains/ip's go. You need to assign a name to your domain aswell.<br>
For more settings read the README.md in the json folder.

## Setup
### Windows
If you are on Windows this script should work without any modification. If you are having trouble/failing ICMP on the pings<br>
then you can try running as admin, other wise jump into the **main.go** and CTRL+F the line `pinger.SetPrivileged(true)` and<br>
comment it out<br>

### Linux
Same as above should apply but if you are having issues with ICMP pings you should jump into the **main.go** and CTRL+F<br>
the line `pinger.SetPrivileged(true)` and comment it out. Then in your linux shell set the following option<br>
```
setcap cap_net_raw=+ep /path/to/your/compiled/binary
```
If you are encountering other problems with the pinging function have a read through the library's installation and supported operating<br>
systems section in the README.md @ [Go-Ping](https://github.com/go-ping/ping)

## TODO
1. Split instance ping and domain ping into different sections so it isn't as messy. i.e, instances.json only contains the addresses for the direct instance.<br>
   I will add a domains.json file where websites/domains can be stored and then the script will do two pings, one on the actual instance servers, and then <br>
   another on the domains/websites. I will then make the email have 2 sections, one for instances and one for domains.
