package task

import (
	"github.com/cloudfoundry-community/go-cfclient"
)

type Variables struct {
	Org   cfclient.Org
	Space cfclient.Space
	App   cfclient.App
}
