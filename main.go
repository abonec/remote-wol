package main

import (
	"net/http"
	"html/template"
	"log"
	"github.com/sabhiram/go-wol"
)

const defaultMac = "48:e2:44:f6:3c:f3"
const bcast = "255.255.255.255:9"

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
	return wol.SendMagicPacket(addr, bcast, "")
}

func failError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	inputTemplate, err := template.New("input_template").Parse(TEMPLATE_STRING)
	failError(err)
	http.HandleFunc("/", indexHandler(inputTemplate))
	err = http.ListenAndServe(":8081", nil)
	failError(err)
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
