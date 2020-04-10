package task

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"time"
)

type Variables struct {
	Org   cfclient.Org
	Space cfclient.Space
	App   cfclient.App
}

func (v Variables) TimeSince(timeString string) time.Duration {
	t, err := v.toTime(timeString)
	if err != nil {
		fmt.Printf("Error parsing time '%s': %s\n", timeString, err)
		return 0
	}
	return time.Since(t)
}

func (v Variables) TimeParseDuration(durString string) time.Duration {
	d, err := time.ParseDuration(durString)
	if err != nil {
		fmt.Printf("Error parsing duration '%s': %s\n", durString, err)
		return 0
	}
	return d
}

func (v Variables) toTime(timeString string) (time.Time, error) {
	// taken from cfclient itself
	possibleFormats := [...]string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05 -0700", "2006-01-02 15:04:05 MST"}

	for _, possibleFormat := range possibleFormats {
		if value, err := time.Parse(possibleFormat, timeString); err == nil {
			return value, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date string: '%s'", timeString)
}
