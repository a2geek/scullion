package task

import (
	"scullion/config"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/cloudfoundry-community/go-cfclient"
)

type Metadata struct {
	Name      string
	Client    *cfclient.Client
	OrgExpr   *vm.Program
	SpaceExpr *vm.Program
	AppExpr   *vm.Program
	Action    func(Item)
}

func NewMetadata(taskDef config.TaskDef, client *cfclient.Client, action func(Item)) (Metadata, error) {
	orgExpr, err := expr.Compile(taskDef.Filters.Organization)
	if err != nil {
		return Metadata{}, err
	}
	spaceExpr, err := expr.Compile(taskDef.Filters.Space)
	if err != nil {
		return Metadata{}, err
	}
	appExpr, err := expr.Compile(taskDef.Filters.Application)
	if err != nil {
		return Metadata{}, err
	}
	metadata := Metadata{
		Name:      taskDef.Name,
		Client:    client,
		OrgExpr:   orgExpr,
		SpaceExpr: spaceExpr,
		AppExpr:   appExpr,
		Action:    action,
	}
	return metadata, nil
}
