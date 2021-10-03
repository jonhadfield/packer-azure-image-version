package packer_azure_image_version

import (
	"github.com/sirupsen/logrus"
	"os"
)

const strAtLeastOnePathRequired = "at least one path is required"

func init() {
	lvl, ok := os.LookupEnv("PAIM_LOG")
	// LOG_LEVEL not set, default to info
	if !ok {
		lvl = "info"
	}

	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}

	logrus.SetLevel(ll)
}
