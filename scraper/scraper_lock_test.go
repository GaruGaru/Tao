package scraper

import (
	"testing"
	"fmt"
	"math/rand"
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
