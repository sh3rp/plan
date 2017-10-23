package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/plan"
)

var dbdir string
var port int
var password string

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.IntVar(&port, "port", 8080, "Port to run plan server on")
	flag.StringVar(&dbdir, "dir", "", "Directory containings the plan database")
	flag.StringVar(&password, "pass", "", "Password to authenticate new entry posting")
	flag.Parse()

	if dbdir == "" {
		fmt.Println("You must specify a -dir paramter with the directory to put the plan database in.")
		return
	}

	if password == "" {
		fmt.Println("You must specify a -pass parameter that is the password used to add new plans.")
		return
	}

	var info *plan.PlanInfo

	if _, err := os.Stat(dbdir + "/plan.db"); os.IsNotExist(err) {
		fmt.Println("")
		fmt.Println("Looks like this is your first time initializing ")
		fmt.Println("your plan database.  Let's set up some information")
		fmt.Println("for your plan!")
		fmt.Println("")

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What is your handle? ")
		handle, _ := reader.ReadString('\n')
		handle = strings.TrimSpace(handle)
		fmt.Printf("What's your real name? ")
		realName, _ := reader.ReadString('\n')
		realName = strings.TrimSpace(realName)
		fmt.Printf("What's your location? ")
		location, _ := reader.ReadString('\n')
		location = strings.TrimSpace(location)

		info = &plan.PlanInfo{
			Handle:   handle,
			RealName: realName,
			Location: location,
		}

	}

	log.Info().Msg("Starting database")
	db, err := plan.NewBoltPlanDB(dbdir)

	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return
	}

	if info != nil {
		db.SetInfo(info)
	}

	log.Info().Msgf("Starting webserver on port %d", port)
	ws := plan.WebService{
		PlanDB: db,
	}

	ws.Start(port, password)
}
