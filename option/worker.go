package option

type WorkerPools struct {
	OrgPool    int `long:"org-pool" env:"ORG_POOL" default:"1" description:"set the number of organization workers in the pool"`
	SpacePool  int `long:"space-pool" env:"SPACE_POOL" default:"1" description:"set the number of space workers in the pool"`
	AppPool    int `long:"app-pool" env:"APP_POOL" default:"1" description:"set the number of application workers in the pool"`
	ActionPool int `long:"action-pool" env:"ACTION_POOL" default:"1" description:"set the number of action (stop/start) workers in the pool"`
}
