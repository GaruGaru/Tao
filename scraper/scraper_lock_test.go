package scraper

import (
	"testing"
	"fmt"
	"math/rand"
	"github.com/GaruGaru/Tao/tests"
)

func TestFileSystemLockObtainRelease(t *testing.T) {

	lock := FileSystemLock{fmt.Sprintf("test.lock.%d", rand.Int())}

	defer lock.Release()

	err := lock.Obtain()

	if err != nil {
		t.Log("Unable to obtain a new file system lock")
		t.FailNow()
	}

	err = lock.Release()

	if err != nil {
		t.Log("Unable to release file system lock")
		t.FailNow()
	}

}

func TestFileSystemLockObtainFail(t *testing.T) {

	lock := FileSystemLock{fmt.Sprintf("test.lock.%d", rand.Int())}

	defer lock.Release()

	err := lock.Obtain()
	if err != nil {
		t.Log("Unable to obtain a new file system lock")
		t.FailNow()
	}

	err = lock.Obtain()

	if err == nil {
		t.Log("File system lock obtained 2 times")
		t.FailNow()
	}

}

func TestFileSystemLockReleaseFail(t *testing.T) {

	lock := FileSystemLock{fmt.Sprintf("test.lock.%d", rand.Int())}

	defer lock.Release()

	err := lock.Obtain()
	if err != nil {
		t.Log("Unable to obtain a new file system lock")
		t.FailNow()
	}

	err = lock.Release()
	if err != nil {
		t.Log("Unable to release a new file system lock")
		t.FailNow()
	}

	err = lock.Release()

	if err == nil {
		t.Log("File system lock released 2 times")
		t.FailNow()
	}

}

///////////

func TestRedisLockObtainRelease(t *testing.T) {

	lock := RedisDojoScraperLock{
		Redis:   *tests.TestRedisClient(t),
		LockKey: fmt.Sprintf("lock_%d", rand.Int()),
	}

	defer lock.Release()

	err := lock.Obtain()

	if err != nil {
		t.Log("Unable to obtain a new redis lock")
		t.FailNow()
	}

	err = lock.Release()

	if err != nil {
		t.Log("Unable to release redis lock")
		t.FailNow()
	}

}

func TestRedisLockObtainFail(t *testing.T) {

	lock := RedisDojoScraperLock{
		Redis:   *tests.TestRedisClient(t),
		LockKey: fmt.Sprintf("lock_%d", rand.Int()),
	}

	defer lock.Release()

	err := lock.Obtain()
	if err != nil {
		t.Log("Unable to obtain a new redis lock")
		t.FailNow()
	}

	err = lock.Obtain()

	if err == nil {
		t.Log("Redis lock obtained 2 times")
		t.FailNow()
	}

}

func TestRedisLockReleaseFail(t *testing.T) {

	lock := RedisDojoScraperLock{
		Redis:   *tests.TestRedisClient(t),
		LockKey: fmt.Sprintf("lock_%d", rand.Int()),
	}

	defer lock.Release()

	err := lock.Obtain()
	if err != nil {
		t.Log("Unable to obtain a new redis lock")
		t.FailNow()
	}

	err = lock.Release()
	if err != nil {
		t.Log("Unable to release a new redis lock")
		t.FailNow()
	}

	err = lock.Release()

	if err == nil {
		t.Log("Redis lock released 2 times ")
		t.FailNow()
	}

}
