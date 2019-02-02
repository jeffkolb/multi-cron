package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	//"github.com/robfig/cron"
	"gopkg.in/robfig/cron.v2"
)

func execute(command string, args []string) {

	println("Doing Thangs")
	println("executing:", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	cmd.Wait()
}

func create() (cr *cron.Cron, wgr *sync.WaitGroup) {
	var schedule = os.Args[1]
	var command = os.Args[2]
	var args = os.Args[3:len(os.Args)]

	execute(command, args)
	wg := &sync.WaitGroup{}
	c := cron.New()
	println("new cron:", schedule)
	c.AddFunc(schedule, func() {
		wg.Add(1)
		execute(command, args)
		wg.Done()
	})

	return c, wg
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

func main() {
	/*
			All key/value pairs for debugging purposes
			for _, e := range os.Environ() {
				pair := strings.Split(e, "=")
				// pair[0] is key, pair[1] is value
				fmt.Println(pair[0] + "=" + pair[1])
		    }
	*/

	// ToDo: Test Connections to Azure and the target database

	keyVal, shouldRun := os.LookupEnv("CLEAN_AT_START")
	if shouldRun && strings.ToUpper(keyVal) == "TRUE" {
		// Run clean at start
		fmt.Println("Running clean at application start.")
	}
	keyVal, shouldRun = os.LookupEnv("BACKUP_AT_START")
	if shouldRun && strings.ToUpper(keyVal) == "TRUE" {
		// Run backup at start
		println("Running backup at application start.")
	}
	keyVal, shouldRun = os.LookupEnv("EMPTY_KEY")
	if shouldRun && strings.ToUpper(keyVal) == "TRUE" {
		println("Running backup at application start.")
	}
	keyVal, shouldRun = os.LookupEnv("KEY_NOT_PRESENT")
	if shouldRun && strings.ToUpper(keyVal) == "TRUE" {
		println("Running backup at application start.")
	}
	_ = keyVal
	_ = shouldRun

	c, wg := create()
	// - Schedule them tasks!
	// - Write warnings and such on errors? Maybe kill the container?
	go start(c, wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	println(<-ch)

	stop(c, wg)
}
