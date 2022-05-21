package framework

//instnces.json
type Instances struct {
	Name []string `json:"name"`
	Addr []string `json:"addr"`
}

//smtp.json
type Mail struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"`

	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`

	Alternative string `json:"alternative"`
}
