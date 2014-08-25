package nodeify

import (
	"time"
)

type Module struct {
	Name          string
	LatestVersion string
	UpdatedAt     time.Time
}
