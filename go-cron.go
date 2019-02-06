package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/robfig/cron"
	//"gopkg.in/robfig/cron.v2"
)

type cronEntry struct {
	Schedule     string
	Application  string
	Args         []string
	RunAtStartup bool
}

func start(c *cron.Cron, wg *sync.WaitGroup) {

	c.Start()
}

func stop(c *cron.Cron, wg *sync.WaitGroup) {
	println("Stopping")
	c.Stop()
	println("Waiting")
	wg.Wait()
	println("Exiting")
	os.Exit(0)
}

func execute(command *cronEntry) {

	/*
		t := time.Now().Format(time.RFC1123)
		fmt.Printf("Executing %v%v at %v\n", command.Application,
		strings.Join(command.Args, " "), t)
	*/
	// ToDo: Catch error starting
	cmd := exec.Command(command.Application, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	//ToDo: Catch error stopping
	cmd.Wait()
}

func getCronEntries() (entries []cronEntry) {
	cronEntries := make([]cronEntry, 0)
	runAtStartup := 0
	// Get all cron schedules defined in environmental variables.
	for _, e := range os.Environ() {
		var entry cronEntry
		sch := strings.Split(e, "=")[1]
		if strings.HasPrefix(e, "CRON_SCH") {
			cmdKey := strings.Replace(strings.Split(e, "=")[0], "_SCH_", "_CMD_", 1)
			argsKey := strings.Replace(strings.Split(e, "=")[0], "_SCH_", "_ARGS_", 1)
			cmdVal, cmdExists := os.LookupEnv(cmdKey)
			if cmdExists {
				//Ones starting with '!' are run at startup
				if strings.HasPrefix(sch, "!") {
					entry.RunAtStartup = true
					runAtStartup++
				}
				entry.Schedule = strings.TrimLeft(sch, "!")
				entry.Application = cmdVal
				argsVal, argsExist := os.LookupEnv(argsKey)
				if argsExist {
					entry.Args = strings.Split(argsVal, " ")
				}
				cronEntries = append(cronEntries, entry)
			} else {
				fmt.Printf("No command exists for %v\n", strings.Split(e, "=")[0])
			}
		}
		//_ = entry
	}

	fmt.Printf("%v cron schedules declared\n", len(cronEntries))
	fmt.Printf("%v commands are configured to run at startup. \n", runAtStartup)
	return cronEntries
}

func main() {
	//wg := &sync.WaitGroup{}

	//c.AddFunc("@every 10s", func() { fmt.Println("Every 10s") })
	//c.AddFunc("@every 15s", func() { fmt.Println("Every 15s") })
	// ToDo: Error Handling like WHOA
	cronEntries := getCronEntries()
	for _, e := range cronEntries {
		var instance cronEntry
		instance.Schedule = e.Schedule
		instance.Application = e.Application
		instance.Args = e.Args
		instance.RunAtStartup = e.RunAtStartup

		if instance.RunAtStartup {
			fmt.Printf("running \"%v %v\" at startup\n", instance.Application,
				strings.Join(instance.Args, " "))
			execute(&instance)
		}

		c := cron.New()
		c.AddFunc(e.Schedule, func() { execute(&instance) })
		c.Start()
		_ = c
		_ = e
	}

	// ToDo: the cronEntries slice should return in the numerical order
	// but we should sort the slice by number just in case.
	// Note: No one should ever have loads of cron entries because that wouldn't be "Dockeresque".

	// ToDo: Write warnings and such on errors? Maybe kill the container?

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	println(<-ch)
	// ToDo: Write a graceful termination
	// stop(c, wg)
}
