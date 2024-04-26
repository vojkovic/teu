package daemon

import (
	"fmt"
	"time"
)

func StartDaemon() error {
	// curent time is printed

	fmt.Println("Starting daemon at", time.Now().Format(time.RFC1123))

	ticker := time.NewTicker(20 * time.Second)

	go func() {
		for {
			<-ticker.C
			err := reconcile()
			if err != nil {
				fmt.Println("error reconciling at", time.Now().Format(time.RFC1123), ": ", err)
			}
		}
	}()

	select {}
}