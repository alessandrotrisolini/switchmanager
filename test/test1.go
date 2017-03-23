package main

import (
	"fmt"
	"os"
	"switchmanager/test/testpkg"
	l "switchmanager/logging"
)

func main() {
	
	err := l.LogInit(os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	err = l.LogInit(os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	logg := l.GetLogger()

	logg.Error("CIAO")

	logg.SetDebugLevel(31)

	logg.Trace("test1")
	logg.Debug("test1")
	logg.Info("test1")
	logg.Warn("test1")
	logg.Error("test1")

	testpkg.DoLog()

}
