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
	flag.StringVar(&dbdir, "dir", "", "Directory containing the plan data")
	flag.Parse()

	if dbdir == "" {
		homedir := os.Getenv("HOME")

		dbdir = homedir + "/.plan"

		if _, err := os.Stat(dbdir); os.IsNotExist(err) {
			os.Mkdir(dbdir, 0700)
		}
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

		fmt.Printf("What's your website URL? ")
		website, _ := reader.ReadString('\n')
		website = strings.TrimSpace(website)

		fmt.Printf("What's your avatar URL? ")
		avatar, _ := reader.ReadString('\n')
		avatar = strings.TrimSpace(avatar)

		info = &plan.PlanInfo{
			Handle:    handle,
			RealName:  realName,
			Location:  location,
			Homepage:  website,
			AvatarURL: avatar,
		}

		fmt.Printf("Enter password for posting: ")
		pass, err := gopass.GetPasswd()

		if err != nil {
			fmt.Println("Error capturing password.")
			os.Exit(1)
		}

		hasher := sha1.New()
		hasher.Write(pass)
		passwordBytes := hasher.Sum(nil)
		password = string(passwordBytes)

		err := ioutil.WriteFile(dbdir+"/passwd", passwordBytes, 0600)

		if err != nil {
			fmt.Printf("Error writing password to file: %v\n", err)
			os.Exit(1)
		}
	}

	if password == "" {
		passwordBytes, err := ioutil.ReadFile(dbdir + "/passwd")
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			os.Exit(1)
		}
		password = string(passwordBytes)
	}

	log.Info().Msgf("Plan v%s (api: v%s)", plan.SERVER_VERSION, plan.API_VERSION)

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
