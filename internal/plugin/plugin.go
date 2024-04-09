package plugin

import (
	"fmt"
	"time"
)

func PluginDetectDaemon() {
	for {
		fmt.Printf("%v+\n", time.Now())
		fmt.Printf("nanjing\n")
		time.Sleep(time.Second * 2)
	}
}
