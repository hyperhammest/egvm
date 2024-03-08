package extension

import (
	"time"
)

func Sleep(seconds uint) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func SleepMs(ms uint) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}
