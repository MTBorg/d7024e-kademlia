package refreshtimer

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type RefreshTimer struct {
	bucketIndex int
	restart     chan bool
}

func NewRefreshTimer(bucketIndex int) *RefreshTimer {
	rt := RefreshTimer{}
	rt.bucketIndex = bucketIndex
	rt.restart = make(chan bool)
	return &rt
}

func (rt *RefreshTimer) StartRefreshTimer(doRefresh func(int)) {
	go func() {
		for {
			//t := time.Duration(10) * time.Second // 10s
			//t := time.Minute
			t := time.Hour
			select {
			case <-rt.restart:
				log.Trace().Str("Bucket", fmt.Sprint(rt.bucketIndex)).Msg("Restarted bucket refresh timer")
			case <-time.After(t):
				log.Trace().Str("Bucket", fmt.Sprint(rt.bucketIndex)).Msg("No lookup done in bucket range, refreshing it")
				go doRefresh(rt.bucketIndex)
			}
		}
	}()
}

func (rt *RefreshTimer) RestartRefreshTimer() {
	rt.restart <- true // restart the refresh timer
}
