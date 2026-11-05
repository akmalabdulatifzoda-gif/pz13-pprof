package work

import (
	"log"
	"time"
)

// TimeIt — ручной таймер-декоратор: defer work.TimeIt("name")()
func TimeIt(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %s", name, time.Since(start))
	}
}
