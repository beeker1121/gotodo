package creek

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	timeNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	path := "testdir_" + timeNano

	logger := log.New(New(path+"/test.log", 1), "Creek Test: ", 0)
	logger.Println("Test line 1")
	logger.Println("Test line 2")

	file, err := os.Open(path + "/test.log")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	var fileLogs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLogs = append(fileLogs, scanner.Text())
	}

	if fileLogs[0] != "Creek Test: Test line 1" {
		t.Errorf("Expected first log to be \"Creek Test: Test line 1\", got %s", fileLogs[0])
	}
	if fileLogs[1] != "Creek Test: Test line 2" {
		t.Errorf("Expected first log to be \"Creek Test: Test line 2\", got %s", fileLogs[1])
	}

	if err := os.RemoveAll(path); err != nil {
		t.Error(err)
	}
}

func TestRollover(t *testing.T) {
	timeNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	path := "testdir_" + timeNano
	megabyte := int64(1024 * 1024)

	logger := New(path+"/test.log", 1)

	for i := 0; int64(i) < megabyte; i++ {
		logger.Write([]byte{0x01})
	}

	file, err := os.Open(path + "/test.log")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		t.Error(err)
	}

	if info.Size() != megabyte {
		t.Errorf("Expected file size to be %d, got %d", megabyte, info.Size())
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	if len(files) != 1 {
		t.Errorf("Expected only one file, found %d", len(files))
	}

	logger.Write([]byte{0x02})

	files, err = ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	if len(files) < 2 {
		t.Errorf("Expected to find two or more files, found %d", len(files))
	}

	if err := os.RemoveAll(path); err != nil {
		t.Error(err)
	}
}

func TestGZIP(t *testing.T) {
	timeNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	path := "testdir_" + timeNano

	err := os.MkdirAll(path, 0744)
	if err != nil {
		t.Error(err)
	}

	file, err := os.OpenFile(path+"/test.log", os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		t.Error(err)
	}

	fileString := "This is a test string that will be compressed and decompressed to verify no data loss."

	if _, err = file.Write([]byte(fileString)); err != nil {
		t.Error(err)
	}
	file.Close()

	compressLogFile(path + "/test.log")

	if _, err = os.Stat(path + "/test.log"); !os.IsNotExist(err) {
		t.Errorf("Expected error to return true on IsNotExist check, got false: %s\n", err)
	}

	filegz, err := os.Open(path + "/test.log.gz")
	if err != nil {
		t.Error(err)
	}
	defer filegz.Close()

	gz, err := gzip.NewReader(filegz)
	if err != nil {
		t.Error(err)
	}
	defer gz.Close()

	s, err := ioutil.ReadAll(gz)
	if err != nil {
		t.Error(err)
	}

	if string(s) != fileString {
		t.Error("Compressed file string does not match original string")
	}

	if err := os.RemoveAll(path); err != nil {
		t.Error(err)
	}
}
