# Instance Pinger
This is a ip/domain pinger made in GoLang that sends an alert email on failed pings.<br>
I primarily use this for keeping track of my AWS Lightsail instances as they go down often.<br>

## Usage
In the **json** folder is where your settings are located.<br>
The **smtp.json** contains the smtp settings. Put your mail providers settigns in here.<br>
The **instances.json** is where your domains/ip's go. You need to assign a name to your domain aswell.<br>
For more settings read the README.md in the json folder.
