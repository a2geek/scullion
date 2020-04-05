package task

import "github.com/cloudfoundry-community/go-cfclient"

type TaskVariables struct {
	Org   cfclient.Org
	Space cfclient.Space
	App   cfclient.App
}
