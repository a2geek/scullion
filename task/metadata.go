package task

import (
	"encoding/json"
	"scullion/config"
	"scullion/util"
	"time"

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
	env, err := toEnv(Variables{})
	if err != nil {
		return Metadata{}, err
	}
	options := []expr.Option{
		expr.Env(env),
		expr.AsBool(),

		// Operators override for date comprising.
		expr.Operator("==", "Equal"),
		expr.Operator("<", "Before"),
		expr.Operator("<=", "BeforeOrEqual"),
		expr.Operator(">", "After"),
		expr.Operator(">=", "AfterOrEqual"),

		// Time and duration manipulation.
		expr.Operator("+", "Add"),
		expr.Operator("-", "Sub"),

		// Operators override for duration comprising.
		expr.Operator("==", "EqualDuration"),
		expr.Operator("<", "BeforeDuration"),
		expr.Operator("<=", "BeforeOrEqualDuration"),
		expr.Operator(">", "AfterDuration"),
		expr.Operator(">=", "AfterOrEqualDuration"),
	}

	orgExpr, err := expr.Compile(taskDef.Filters.Organization, options...)
	if err != nil {
		return Metadata{}, err
	}
	spaceExpr, err := expr.Compile(taskDef.Filters.Space, options...)
	if err != nil {
		return Metadata{}, err
	}
	appExpr, err := expr.Compile(taskDef.Filters.Application, options...)
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

func (m *Metadata) IsOrgMatch(vars Variables) (bool, error) {
	return isMatch(m.OrgExpr, vars)
}
func (m *Metadata) IsSpaceMatch(vars Variables) (bool, error) {
	return isMatch(m.SpaceExpr, vars)
}
func (m *Metadata) IsAppMatch(vars Variables) (bool, error) {
	return isMatch(m.AppExpr, vars)
}

func isMatch(pgm *vm.Program, vars Variables) (bool, error) {
	env, err := toEnv(vars)
	if err != nil {
		return false, err
	}
	result, err := expr.Run(pgm, env)
	if err != nil {
		return false, err
	}
	return util.IsTrue(result)
}

type RunEnv struct {
	Org   map[string]interface{}
	Space map[string]interface{}
	App   map[string]interface{}
	datetime
}

func toEnv(vars Variables) (interface{}, error) {
	orgMap, err := toMap(vars.Org)
	if err != nil {
		return nil, err
	}
	spaceMap, err := toMap(vars.Space)
	if err != nil {
		return nil, err
	}
	appMap, err := toMap(vars.App)
	if err != nil {
		return nil, err
	}
	env := RunEnv{
		Org:   orgMap,
		Space: spaceMap,
		App:   appMap,
	}
	return env, nil
}
func toMap(obj interface{}) (map[string]interface{}, error) {
	theJSON, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	theMap := map[string]interface{}{}
	err = json.Unmarshal(theJSON, &theMap)
	if err != nil {
		return nil, err
	}
	return theMap, err
}

// See: https://github.com/antonmedv/expr/blob/master/docs/examples/dates_test.go
type datetime struct{}

func (datetime) Date(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
func (datetime) Duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}
func (datetime) Now() time.Time                                { return time.Now() }
func (datetime) Equal(a, b time.Time) bool                     { return a.Equal(b) }
func (datetime) Before(a, b time.Time) bool                    { return a.Before(b) }
func (datetime) BeforeOrEqual(a, b time.Time) bool             { return a.Before(b) || a.Equal(b) }
func (datetime) After(a, b time.Time) bool                     { return a.After(b) }
func (datetime) AfterOrEqual(a, b time.Time) bool              { return a.After(b) || a.Equal(b) }
func (datetime) Add(a time.Time, b time.Duration) time.Time    { return a.Add(b) }
func (datetime) Sub(a, b time.Time) time.Duration              { return a.Sub(b) }
func (datetime) EqualDuration(a, b time.Duration) bool         { return a == b }
func (datetime) BeforeDuration(a, b time.Duration) bool        { return a < b }
func (datetime) BeforeOrEqualDuration(a, b time.Duration) bool { return a <= b }
func (datetime) AfterDuration(a, b time.Duration) bool         { return a > b }
func (datetime) AfterOrEqualDuration(a, b time.Duration) bool  { return a >= b }
