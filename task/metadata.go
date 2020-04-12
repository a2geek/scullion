package task

import (
	"encoding/json"
	"fmt"
	"scullion/config"
	"scullion/log"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/lxc/lxd/shared/logger"
)

type Metadata struct {
	Name      string
	Client    *cfclient.Client
	OrgExpr   *vm.Program
	SpaceExpr *vm.Program
	AppExpr   *vm.Program
	Action    func(Item)
	Logger    log.Logger
}

func NewMetadata(taskDef config.TaskDef, client *cfclient.Client, action func(Item), logLevel string) (Metadata, error) {
	logger, err := log.NewLogger(taskDef.Name, logLevel)
	if err != nil {
		return Metadata{}, err
	}

	orgExpr, spaceExpr, appExpr, err := compile(taskDef.Filters, logger)
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
		Logger:    logger,
	}
	return metadata, nil
}

func compile(filters config.Filter, logger log.Logger) (*vm.Program, *vm.Program, *vm.Program, error) {
	env, err := toEnv(Variables{}, logger)
	if err != nil {
		return nil, nil, nil, err
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

	orgExpr, err := expr.Compile(filters.Organization, options...)
	if err != nil {
		return nil, nil, nil, err
	}
	spaceExpr, err := expr.Compile(filters.Space, options...)
	if err != nil {
		return nil, nil, nil, err
	}
	appExpr, err := expr.Compile(filters.Application, options...)
	if err != nil {
		return nil, nil, nil, err
	}
	return orgExpr, spaceExpr, appExpr, nil
}

func (m *Metadata) IsOrgMatch(vars Variables) (bool, error) {
	return m.isMatch(m.OrgExpr, vars)
}
func (m *Metadata) IsSpaceMatch(vars Variables) (bool, error) {
	return m.isMatch(m.SpaceExpr, vars)
}
func (m *Metadata) IsAppMatch(vars Variables) (bool, error) {
	return m.isMatch(m.AppExpr, vars)
}

func (m *Metadata) isMatch(pgm *vm.Program, vars Variables) (bool, error) {
	env, err := toEnv(vars, m.Logger)
	if err != nil {
		return false, err
	}
	result, err := expr.Run(pgm, env)
	if err != nil {
		return false, err
	}
	return isTrue(result)
}

func isTrue(i interface{}) (bool, error) {
	switch t := i.(type) {
	case int:
		return i != 0, nil
	case string:
		return i != "", nil
	case bool:
		return i.(bool), nil
	default:
		return false, fmt.Errorf("unable to test type '%s', value '%s'", t, i)
	}
}

type RunEnv struct {
	Org   map[string]interface{}
	Space map[string]interface{}
	App   map[string]interface{}
	datetime
}

func toEnv(vars Variables, logger log.Logger) (interface{}, error) {
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
	env.logger = logger
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
type datetime struct {
	logger log.Logger
}

func (datetime) Date(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		logger.Errorf("date parse: %v", err)
	}
	return t
}
func (datetime) Duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		logger.Errorf("duration parse: %v", err)
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
func (d datetime) Since(s string) time.Duration                { return time.Since(d.Date(s)) }
