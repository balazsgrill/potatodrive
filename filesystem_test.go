package projfero_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/balazsgrill/projfero"
	"github.com/spf13/afero"
)

type testInstance struct {
	t         *testing.T
	location  string
	fs        afero.Fs
	closer    io.Closer
	closechan chan bool
}

func newTestInstance(t *testing.T) *testInstance {
	location := t.TempDir()
	os.RemoveAll(location)
	os.MkdirAll(location, 0x777)
	return &testInstance{
		t:         t,
		location:  location,
		fs:        afero.NewMemMapFs(),
		closechan: make(chan bool),
	}
}

func (i *testInstance) start() {
	started := make(chan bool)
	var err error
	go func() {
		i.closer, err = projfero.StartProjecting(i.location, i.fs)
		started <- true
		<-i.closechan
		i.closer.Close()
	}()
	<-started
	if err != nil {
		log.Fatal(err)
	}
}

func (i *testInstance) osWriteFile(filename string, content string) error {
	return exec.Command("cmd", "/c", "echo", content, ">", i.location+"\\"+filename).Run()
}

func (i *testInstance) osRemoveFile(filename string) error {
	return exec.Command("cmd", "/c", "del", i.location+"\\"+filename).Run()
}

func (i *testInstance) osCreateDir(filename string) error {
	return exec.Command("cmd", "/c", "mkdir", i.location+"\\"+filename).Run()
}

func (i *testInstance) osRemoveDir(filename string) error {
	return exec.Command("cmd", "/c", "rmdir", i.location+"\\"+filename).Run()
}

func (i *testInstance) stop() {
	i.closechan <- true
}

func TestExistingFileOnBackend(t *testing.T) {
	instance := newTestInstance(t)

	data := []byte("something")
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, data, 0x777)
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	defer instance.stop()

	data2, err := os.ReadFile(instance.location + "\\" + filename)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, data2) {
		t.Errorf("expected %v, got %v", data, data2)
	}
}

func TestFileCreation(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()
	defer instance.stop()

	filename := "test.txt"
	data := "something"
	err := instance.osWriteFile(filename, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := afero.ReadFile(instance.fs, filename)
	if err != nil {
		t.Fatal(err)
	}

	if data != strings.TrimSpace(string(data2)) {
		t.Errorf("expected '%s', got '%s'", data, string(data2))
	}
}

func TestUpdateExistingFileOnBackend(t *testing.T) {
	instance := newTestInstance(t)

	data := "something"
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, []byte(data), 0x777)
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	defer instance.stop()

	data = "somethingelse"
	err = instance.osWriteFile(filename, data)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := afero.ReadFile(instance.fs, filename)
	if err != nil {
		t.Fatal(err)
	}

	if data != strings.TrimSpace(string(data2)) {
		t.Errorf("expected %s, got %s", data, string(data2))
	}
}

func TestDeleteExistingFileOnBackend(t *testing.T) {
	instance := newTestInstance(t)
	data := "something"
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, []byte(data), 0x777)
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	defer instance.stop()

	err = instance.osRemoveFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	_, err = instance.fs.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			//ok
			return
		}
		t.Fatal(err)
	} else {
		t.Error("File exists")
	}
}

func TestListFiles(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()
	defer instance.stop()

	data := "something"
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, []byte(data), 0x777)
	if err != nil {
		t.Fatal(err)
	}

	filename2 := "test2.txt"
	err = instance.osWriteFile(filename2, data)
	if err != nil {
		t.Fatal(err)
	}

	expected := make(map[string]bool)
	expected[filename] = true
	expected[filename2] = true

	entries, err := os.ReadDir(instance.location)
	if err != nil {
		t.Fatal(err)
	}

	actual := make(map[string]bool)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		actual[entry.Name()] = true
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
