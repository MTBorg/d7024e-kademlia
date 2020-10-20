package refreshtimer

import (
	"fmt"
	"os"
	"strconv"
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
		// get the refresh time (in seconds) from env or default to 1 hour
		bucketRefreshTime, err := strconv.Atoi(os.Getenv("BUCKET_REFRESH_TIME"))
		if err != nil {
			log.Error().Msgf("Failed to convert env variable REFRESH_TIME from string to int: %s", err)
			bucketRefreshTime = 3600
		}
		for {
			//t := time.Duration(10) * time.Second // 10s
			//t := time.Minute
			t := time.Duration(bucketRefreshTime) * time.Second
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
