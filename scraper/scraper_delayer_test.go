package scraper

import (
	"fmt"
	"github.com/GaruGaru/Tao/tests"
	"math/rand"
	"testing"
	"time"
)

func TestScraperDelayerRedis(t *testing.T) {

	redisClient := *tests.TestRedisClient(t)
	timeKey := fmt.Sprintf("delay_%d", rand.Int())

	defer redisClient.Del(timeKey)

	delayer := RedisScraperDelayer{
		Redis:   redisClient,
		Delay:   2 * time.Second,
		TimeKey: timeKey,
	}

	canRun, err := delayer.CanRun()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if !canRun {
		t.Log("Can't run is false on first run.")
		t.FailNow()
	}

	delayer.Refresh() // Delay

	canRun, err = delayer.CanRun()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if canRun {
		t.Log("Can run but should be false")
		t.FailNow()
	}

	delayer.Refresh() // Delay

	time.Sleep(3 * time.Second)
	canRun, err = delayer.CanRun()

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if !canRun {
		t.Log("Can't run but should be false")
		t.FailNow()
	}
}
