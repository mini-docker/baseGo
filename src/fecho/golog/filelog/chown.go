// +build !linux

package filelog

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
