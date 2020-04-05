package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jessevdk/go-flags"
)

func main() {
	tasks := make([]Task, 0)
	options := NewOptions(&tasks)
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	_, err := parser.Parse()
	if err != nil {
		if !flags.WroteHelp(err) {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("Please supply a file or an environment variable for configuration")
		os.Exit(1)
	}
	if options.TaskOptions.Validate {
		if !Validate(tasks) {
			fmt.Println("Some tasks and schedules DID NOT pass validation.")
			os.Exit(1)
		} else {
			fmt.Println("All tasks and schedules passed validation.")
			os.Exit(0)
		}
	}

	orgChan := make(chan TaskItem)
	spaceChan := make(chan TaskItem)
	appChan := make(chan TaskItem)
	actionChan := make(chan TaskItem)

	for i := 0; i < options.OrgPool; i++ {
		go orgWorker(i, orgChan, spaceChan)
	}
	for i := 0; i < options.SpacePool; i++ {
		go spaceWorker(i, spaceChan, appChan)
	}
	for i := 0; i < options.AppPool; i++ {
		go appWorker(i, appChan, actionChan)
	}
	for i := 0; i < options.ActionPool; i++ {
		go actionWorker(i, actionChan)
	}

	c := &cfclient.Config{
		ApiAddress:        options.API,
		Username:          options.Username,
		Password:          options.Password,
		SkipSslValidation: options.SkipSslValidation,
	}
	client, err := cfclient.NewClient(c)
	if err != nil {
		panic(err)
	}

	for i, task := range tasks {
		go taskWorker(i, task, client, orgChan)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}

func createMetadata(task Task, client *cfclient.Client) (TaskMetadata, error) {
	orgExpr, err := expr.Compile(task.Filters.Organization)
	if err != nil {
		return TaskMetadata{}, err
	}
	spaceExpr, err := expr.Compile(task.Filters.Space)
	if err != nil {
		return TaskMetadata{}, err
	}
	appExpr, err := expr.Compile(task.Filters.Application)
	if err != nil {
		return TaskMetadata{}, err
	}
	metadata := TaskMetadata{
		Name:      task.Name,
		Client:    client,
		OrgExpr:   orgExpr,
		SpaceExpr: spaceExpr,
		AppExpr:   appExpr,
		Action:    logAction,
	}
	return metadata, nil
}

func taskWorker(num int, task Task, client *cfclient.Client, orgChan chan<- TaskItem) {
	fmt.Printf("Launched task worker %d\n", num)
	metadata, err := createMetadata(task, client)
	if err != nil {
		panic(err)
	}
	dur, err := time.ParseDuration(task.Schedule.Frequency)
	if err != nil {
		panic(err)
	}
	for t := range time.Tick(dur) {
		fmt.Printf("[%s] Tick at %s\n", task.Name, t)
		taskItem := TaskItem{
			Metadata: metadata,
		}
		orgChan <- taskItem
	}
}

func logAction(taskItem TaskItem) {
	fmt.Println(taskItem.Variables)
}

type TaskItem struct {
	Variables TaskVariables
	Metadata  TaskMetadata
}
type TaskVariables struct {
	Org   cfclient.Org
	Space cfclient.Space
	App   cfclient.App
}
type TaskMetadata struct {
	Name      string
	Client    *cfclient.Client
	OrgExpr   *vm.Program
	SpaceExpr *vm.Program
	AppExpr   *vm.Program
	Action    func(TaskItem)
}

func IsTrue(i interface{}) (bool, error) {
	switch t := i.(type) {
	case int:
		return i != 0, nil
	case string:
		return i != "", nil
	case bool:
		return i.(bool), nil
	default:
		return false, fmt.Errorf("unable to test type '%s'", t)
	}
}
func isTrue(i interface{}) bool {
	b, err := IsTrue(i)
	if err != nil {
		panic(err)
	}
	return b
}

func orgWorker(num int, orgChan <-chan TaskItem, spaceChan chan<- TaskItem) {
	fmt.Printf("Launched org worker %d\n", num)
	for {
		taskItem := <-orgChan
		orgs, err := taskItem.Metadata.Client.ListOrgs()
		if err != nil {
			panic(err)
		}
		for _, org := range orgs {
			variables := TaskVariables{
				Org: org,
			}
			result, err := expr.Run(taskItem.Metadata.OrgExpr, variables)
			if err != nil {
				panic(err)
			}
			if isTrue(result) {
				fmt.Printf("[%s] Matched org '%s'\n", taskItem.Metadata.Name, org.Name)
				newTask := TaskItem{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				spaceChan <- newTask
			} else {
				fmt.Printf("[%s] Skipping org '%s'\n", taskItem.Metadata.Name, org.Name)
			}
		}
	}
}
func spaceWorker(num int, spaceChan <-chan TaskItem, appChan chan<- TaskItem) {
	fmt.Printf("Launched space worker %d\n", num)
	for {
		taskItem := <-spaceChan
		q := url.Values{
			"organization_guid": []string{taskItem.Variables.Org.Guid},
		}
		spaces, err := taskItem.Metadata.Client.ListSpacesByQuery(q)
		if err != nil {
			panic(err)
		}
		for _, space := range spaces {
			variables := TaskVariables{
				Org:   taskItem.Variables.Org,
				Space: space,
			}
			result, err := expr.Run(taskItem.Metadata.SpaceExpr, variables)
			if err != nil {
				panic(err)
			}
			if isTrue(result) {
				fmt.Printf("[%s] Matched space '%s' in org '%s'\n", taskItem.Metadata.Name, space.Name, variables.Org.Name)
				newTask := TaskItem{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				appChan <- newTask
			} else {
				fmt.Printf("[%s] Skipping space '%s' in org '%s'\n", taskItem.Metadata.Name, space.Name, variables.Org.Name)
			}
		}
	}
}
func appWorker(num int, appChan <-chan TaskItem, actionChan chan<- TaskItem) {
	fmt.Printf("Launched app worker %d\n", num)
	for {
		taskItem := <-appChan
		q := url.Values{
			"space_guid": []string{taskItem.Variables.Space.Guid},
		}
		apps, err := taskItem.Metadata.Client.ListAppsByQuery(q)
		if err != nil {
			panic(err)
		}
		for _, app := range apps {
			variables := TaskVariables{
				Org:   taskItem.Variables.Org,
				Space: taskItem.Variables.Space,
				App:   app,
			}
			result, err := expr.Run(taskItem.Metadata.AppExpr, variables)
			if err != nil {
				panic(err)
			}
			if isTrue(result) {
				fmt.Printf("[%s] Matched app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name, app.Name,
					variables.Space.Name, variables.Org.Name)
				newTask := TaskItem{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				actionChan <- newTask
			} else {
				fmt.Printf("[%s] Skipped app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name, app.Name,
					variables.Space.Name, variables.Org.Name)
			}
		}
	}
}
func actionWorker(num int, actionChan <-chan TaskItem) {
	fmt.Printf("Launched action worker %d\n", num)
	for {
		taskItem := <-actionChan
		fmt.Printf("[%s] Performing action on app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name,
			taskItem.Variables.App.Name, taskItem.Variables.Space.Name, taskItem.Variables.Org.Name)
		taskItem.Metadata.Action(taskItem)
	}
}
