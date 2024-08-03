package filesystem_test

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/balazsgrill/potatodrive/filesystem"
	"github.com/spf13/afero"
)

type testInstance struct {
	t         *testing.T
	location  string
	fs        afero.Fs
	closer    filesystem.Virtualization
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
		i.closer, err = filesystem.StartProjecting(i.location, i.fs)
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

func TestExistingFolderOnBackend(t *testing.T) {
	instance := newTestInstance(t)

	foldername := "test"
	instance.fs.Mkdir(foldername, 0x777)

	instance.start()
	defer instance.stop()

	stat, err := os.Stat(instance.location + "\\" + foldername)
	if err != nil {
		t.Fatal(err)
	}

	if stat.IsDir() != true {
		t.Error("Not a directory")
	}
}

func TestFolderCreation(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()
	defer instance.stop()

	foldername := "test"
	err := instance.osCreateDir(foldername)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := instance.fs.Stat(foldername)
	if err != nil {
		t.Fatal(err)
	}
	if stat.IsDir() != true {
		t.Error("Not a directory")
	}
}

func TestCreatedOnBackend(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()
	defer instance.stop()

	data := []byte("something")
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, data, 0x777)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := os.ReadFile(instance.location + "\\" + filename)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, data2) {
		t.Errorf("expected %v, got %v", data, data2)
	}

}

func TestChangedOnBackend(t *testing.T) {
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

	// sleep for a bit to ensure that the file timestamp is different
	time.Sleep(time.Second)
	data = []byte("somethingelse")
	err = afero.WriteFile(instance.fs, filename, data, 0x777)
	if err != nil {
		t.Fatal(err)
	}

	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	data2, err = os.ReadFile(instance.location + "\\" + filename)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, data2) {
		t.Errorf("expected %v, got %v", data, data2)
	}
}

func TestDeletedOnBackend(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()
	defer instance.stop()
	data := []byte("something")
	filename := "test.txt"
	err := afero.WriteFile(instance.fs, filename, data, 0x777)
	if err != nil {
		t.Fatal(err)
	}
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

func TestUpdatedLocallyWhileOffline(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()

	data := []byte("something")
	filename := "test.txt"
	err := instance.osWriteFile(filename, string(data))

	if err != nil {
		t.Fatal(err)
	}

	instance.stop()
	time.Sleep(time.Second)

	data = []byte("somethingelse")
	err = instance.osWriteFile(filename, string(data))
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	data2, err := afero.ReadFile(instance.fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != strings.TrimSpace(string(data2)) {
		t.Errorf("expected '%s', got '%s'", string(data), string(data2))
	}
	instance.stop()
}

func TestRemoveFolder(t *testing.T) {
	foldername := "test"
	instance := newTestInstance(t)
	err := instance.fs.Mkdir(foldername, 0x777)
	if err != nil {
		t.Fatal(err)
	}
	instance.start()
	defer instance.stop()

	file, err := os.Stat(instance.location + "\\" + foldername)
	if err != nil {
		t.Fatal(err)
	}
	if file.IsDir() != true {
		t.Error("Not a directory")
	}

	err = instance.osRemoveDir(foldername)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.fs.Stat(foldername)
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

func TestDeletedOnBackendWhileOffline(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()

	data := []byte("something")
	filename := "test.txt"
	err := instance.osWriteFile(filename, string(data))

	if err != nil {
		t.Fatal(err)
	}

	instance.stop()
	time.Sleep(time.Second)

	err = instance.fs.Remove(filename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.fs.Stat(filename)
	if !os.IsNotExist(err) {
		t.Error("remote file exists")
	}

	instance.start()
	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(instance.location + "\\" + filename)
	if !os.IsNotExist(err) {
		t.Error("local file exists")
	}
	_, err = instance.fs.Stat(filename)
	if !os.IsNotExist(err) {
		t.Error("remote file exists")
	}

	instance.stop()
}

func TestDeletedLocallyWhileOffline(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()

	data := []byte("something")
	filename := "test.txt"
	err := instance.osWriteFile(filename, string(data))

	if err != nil {
		t.Fatal(err)
	}

	instance.stop()
	time.Sleep(time.Second)

	err = instance.osRemoveFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(instance.location + "\\" + filename)
	if os.IsNotExist(err) {
		t.Error("File should be restored locally")
	}
	_, err = instance.fs.Stat(filename)
	if os.IsNotExist(err) {
		t.Error("remote file should not be removed")
	}

	instance.stop()
}

func TestConflictWhileOfflineLocalNewer(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()

	data := []byte("something")
	filename := "test.txt"
	err := instance.osWriteFile(filename, string(data))

	if err != nil {
		t.Fatal(err)
	}

	instance.stop()
	time.Sleep(time.Second)

	data2 := []byte("something2")
	data3 := []byte("something3")

	err = afero.WriteFile(instance.fs, filename, data3, 0x777)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	err = instance.osWriteFile(filename, string(data2))
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	data4, err := afero.ReadFile(instance.fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data2) != strings.TrimSpace(string(data4)) {
		t.Errorf("expected '%s', got '%s'", string(data2), string(data4))
	}
	data5, err := os.ReadFile(instance.location + "\\" + filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data2) != strings.TrimSpace(string(data5)) {
		t.Errorf("expected '%s', got '%s'", string(data2), string(data5))
	}

	instance.stop()
}

func TestConflictWhileOfflineRemoteNewer(t *testing.T) {
	instance := newTestInstance(t)
	instance.start()

	data := []byte("something")
	filename := "test.txt"
	err := instance.osWriteFile(filename, string(data))

	if err != nil {
		t.Fatal(err)
	}

	instance.stop()
	time.Sleep(time.Second)

	data2 := []byte("something2")
	data3 := []byte("something3")

	err = instance.osWriteFile(filename, string(data2))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	err = afero.WriteFile(instance.fs, filename, data3, 0x777)
	if err != nil {
		t.Fatal(err)
	}

	instance.start()
	err = instance.closer.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	data4, err := afero.ReadFile(instance.fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data3) != strings.TrimSpace(string(data4)) {
		t.Errorf("expected '%s', got '%s'", string(data3), string(data4))
	}
	data5, err := os.ReadFile(instance.location + "\\" + filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data3) != strings.TrimSpace(string(data5)) {
		t.Errorf("expected '%s', got '%s'", string(data3), string(data5))
	}

	instance.stop()
}
