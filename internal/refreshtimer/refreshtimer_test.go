package refreshtimer_test

import (
	"kademlia/internal/refreshtimer"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DoerMock struct {
	mock.Mock
}

func (m *DoerMock) Do(bucketIndex int) {
	m.Called(bucketIndex)
	return
}

func TestNewRefreshTimer(t *testing.T) {
	bucketIndex := 20
	rt := refreshtimer.NewRefreshTimer(bucketIndex)

	// should return a refresh timer
	assert.IsType(t, refreshtimer.RefreshTimer{}, *rt)
}

func TestStartRefreshTimer(t *testing.T) {
	var rt *refreshtimer.RefreshTimer

	// should not fail if the refresh time env var is not set
	rt = refreshtimer.NewRefreshTimer(10)
	rt.StartRefreshTimer(func(int) { return })
	time.Sleep(time.Second)

	// should refresh the bucket (by calling the do function) if the timer is
	// not restarted
	doerMock := new(DoerMock)
	rt = refreshtimer.NewRefreshTimer(20)
	os.Setenv("BUCKET_REFRESH_TIME", "1")
	doerMock.On("Do", 20).Return()
	rt.StartRefreshTimer(func(b int) { doerMock.Do(b) })
	// sleep to trigger refresh
	time.Sleep(time.Second * 2)
	doerMock.AssertExpectations(t)
	doerMock.AssertNumberOfCalls(t, "Do", 1)

	// should be able to restart the refresh timer
	doerMock2 := new(DoerMock)
	rt = refreshtimer.NewRefreshTimer(30)
	doerMock2.On("Do", 30).Return()
	os.Setenv("BUCKET_REFRESH_TIME", "1")
	rt.StartRefreshTimer(func(b int) { doerMock2.Do(b) })
	time.Sleep(time.Duration(time.Millisecond) * 500)
	rt.RestartRefreshTimer()
	time.Sleep(time.Duration(time.Millisecond) * 700)
	doerMock2.AssertNotCalled(t, "Do")
	time.Sleep(time.Duration(time.Millisecond) * 500)
	doerMock2.AssertExpectations(t)
	doerMock2.AssertNumberOfCalls(t, "Do", 1)
}
