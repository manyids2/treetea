package task

import (
	"log"
	"os"
)

// OpenLog NOTE: Remember to `defer f.Close()`
func OpenLog(path string) (f *os.File) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
