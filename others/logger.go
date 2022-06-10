package others

import (
	"io"
	"log"
	"os"
	"time"
)

func SetupLogger() error {
	if !TheConfig.Log {
		return nil
	}

	// bind log file
	f, err := os.OpenFile(ConfigDir+"logs/"+time.Now().Format("20060201_15_04_05")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	wrt := io.MultiWriter(os.Stderr, f)
	log.SetOutput(wrt)

	return nil
}
