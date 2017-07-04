package main

import (
	"net/http"
	"html/template"
	"log"
	"github.com/sabhiram/go-wol"
	"github.com/abonec/go-ping"
	"time"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultMac = "f4:f2:6d:02:3f:7e"
	defaultIp  = "192.168.255.14"
	bcast      = "255.255.255.255:9"
)

type Page struct {
	MacAddr string
	Message string
}

func indexHandler(tmp *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmp.Execute(w, Page{MacAddr: defaultMac})
		} else if r.Method == "POST" {
			addr := r.FormValue("mac_addr")
			message := "WOL sent to: " + addr
			err := sendMagicPacket(addr)
			if err != nil {
				message = err.Error()
			}
			tmp.Execute(w, Page{MacAddr: defaultMac, Message: message})
		}
	}
}
func sendMagicPacket(addr string) error {
	return wol.SendMagicPacket(addr, bcast, "eth0")
}

func failError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func turnPowerOn() {
	turnRelay()
	//sendMagicPacket(defaultMac)
}

func main() {
	go handleSignals()
	go startTelegramBot()
	inputTemplate, err := template.New("input_template").Parse(TEMPLATE_STRING)
	failError(err)
	http.HandleFunc("/", indexHandler(inputTemplate))
	err = http.ListenAndServe(":8081", nil)
	failError(err)
}

func pingMachine() bool {
	pinger, err := ping.NewPinger(defaultIp)
	failError(err)
	pinger.Count = 1
	pinger.Timeout = time.Second
	pinger.Run()
	if pinger.Statistics().PacketsRecv > 0 {
		return true
	} else {
		return false
	}
}

func handleSignals() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	for sig := range sigs {
		log.Println(sig)
		haltTelegramBot()
		os.Exit(1)
	}
}

const TEMPLATE_STRING = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width">
</head>
<body>
<form method='post' >
  <div>
  	{{.Message}}
  </div>
  <div>
	  <label for= 'mac_addr'>mac addr:</label>
	  <input id ='mac_addr' name='mac_addr' value='{{.MacAddr}}'/>
  </div>
  <div>
	  <input type='submit'/>
  </div>
  </form>
</body>
</html>
`
