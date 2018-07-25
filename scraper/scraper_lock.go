package scraper

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

type DojoScraperLock interface {
	Obtain() error
	Release() error
}

type FileSystemLock struct {
	LockFile string
}

func (l FileSystemLock) Obtain() error {
	if _, err := os.Stat(l.LockFile); err == nil {
		return fmt.Errorf("locked")
	} else {
		_, err := os.Create(l.LockFile)
		return err
	}
}

func (l FileSystemLock) Release() error {
	return os.Remove(l.LockFile)
}

type RedisDojoScraperLock struct {
	Redis   redis.Client
	LockKey string
}

func (l RedisDojoScraperLock) Obtain() error {
	res := l.Redis.SetNX(l.LockKey, "lock", 0)
	if res.Err() != nil {
		return res.Err()
	} else if res.Val() == false {
		return fmt.Errorf("unable to obtain redis lock for key: %s", l.LockKey)
	} else {
		return nil
	}
}

func (l RedisDojoScraperLock) Release() error {
	r := l.Redis.Del(l.LockKey)
	if r.Err() != nil {
		return r.Err()
	} else {
		v, e := r.Result()

		if e != nil && e != redis.Nil {
			return e
		} else if v == 1 {
			return nil
		} else {
			return fmt.Errorf("unable to release redis lock for key %s", l.LockKey)
		}

	}
	return r.Err()
}

type TestingLock struct {
}

func (l TestingLock) Obtain() error {
	return nil
}

func (l TestingLock) Release() error {
	return nil
}
