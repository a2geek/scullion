package option

import "github.com/cloudfoundry-community/go-cfclient"

type CloudFoundryOptions struct {
	API               string `short:"a" long:"api" env:"API" description:"API URL"`
	Username          string `short:"u" long:"username" env:"USERNAME" description:"Username"`
	Password          string `short:"p" long:"password" env:"PASSWORD" description:"Password"`
	SkipSslValidation bool   `short:"k" long:"skip-ssl-validation" env:"SKIP_SSL_VALIDATION" description:"Skip SSL validation of Cloud Foundry endpoint; not recommended"`
}

func (cf *CloudFoundryOptions) Client() (*cfclient.Client, error) {
	c := &cfclient.Config{
		ApiAddress:        cf.API,
		Username:          cf.Username,
		Password:          cf.Password,
		SkipSslValidation: cf.SkipSslValidation,
	}
	return cfclient.NewClient(c)
}
