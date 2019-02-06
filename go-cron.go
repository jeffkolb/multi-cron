package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	// "github.com/robfig/cron"
	"gopkg.in/robfig/cron.v2"
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
	cstr := command.Application + " " + strings.Join(command.Args, " ")
	t := time.Now().Format(time.RFC1123)
	fmt.Printf("Running \"%v\" at: %v\n\n", cstr, t)
	// ToDo: Catch error starting
	cmd := exec.Command(command.Application, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	//ToDo: Catch error stopping
	cmd.Wait()

	t = time.Now().Format(time.RFC1123)
	fmt.Printf("\nDone running \"%v\" at: %v\n\n", cstr, t)
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
	}

	fmt.Printf("%v cron schedules declared\n", len(cronEntries))
	fmt.Printf("%v commands are configured to run at startup. \n", runAtStartup)
	return cronEntries
}

func main() {
	fmt.Printf("================= Starting multi-cron =================\n")
	wg := &sync.WaitGroup{}
	// ToDo: Error Handling like WHOA
	c := cron.New()
	cronEntries := getCronEntries()
	wg.Add(1)
	for _, e := range cronEntries {
		var instance cronEntry
		instance.Schedule = e.Schedule
		instance.Application = e.Application
		instance.Args = e.Args
		instance.RunAtStartup = e.RunAtStartup
		cmdStr := instance.Application + " " + strings.Join(instance.Args, " ")
		if instance.RunAtStartup {
			fmt.Printf("Configured to execute \"%v\" at multi-cron start\n", cmdStr)
			execute(&instance)
		}
		f := func() { execute(&instance) }
		c.AddFunc(e.Schedule, f)
	}
	inspect := c.Entries()
	c.Start()
	_ = inspect

	// ToDo: the cronEntries slice should return in the numerical order
	// but we should sort the slice by number just in case.
	// Note: No one should ever have loads of cron entries because that wouldn't be "Dockeresque".

	// ToDo: Write warnings and such on errors? Maybe kill the container?

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	println(<-ch)
	// ToDo: Write a graceful termination
	stop(c, wg)
}
