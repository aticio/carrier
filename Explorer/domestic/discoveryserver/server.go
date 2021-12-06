package explorerserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"carrier/Explorer/domestic/model"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func AddServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	servers := []model.Server{}
	json.NewDecoder(r.Body).Decode(&servers)
	log.Info("Servers came in: ", servers)

	p := []string{}
	for _, server := range servers {
		log.Info("Checking server if exists: ", server.Host, server.Port)
		serverExists := checkSeed(server)

		if !serverExists {
			log.Info(server.Host, server.Port, "does not exists in the db. Initiating adding operation...")
			result := sendServer(server)
			p = append(p, server.Host+":"+server.Port+" - "+result)
		} else {
			p = append(p, server.Host+":"+server.Port+" - already exists")
		}
	}

	var b bytes.Buffer

	for _, res := range p {
		b.WriteString(res)
		b.WriteString("\n")
	}

	fmt.Fprint(w, b.String())
}
