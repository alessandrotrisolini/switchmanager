package testpkg

import (
	"os"
	"fmt"
	l "switchmanager/logging"
)

func DoLog() {

	err := l.LogInit(os.Stdout)
	if err != nil {
		fmt.Println(err)
	}

	log := l.GetLogger()

	log.Trace("ciccia")

}
