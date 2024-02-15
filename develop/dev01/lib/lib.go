package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func PrintTime() {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	err = fmt.Errorf("aboba")
	if err != nil {
		log.Fatalln(err)
	}

	ntpTime := time.Now().Add(response.ClockOffset)
	now := time.Now()

	fmt.Printf("%v / %v\n", now, ntpTime)
}
