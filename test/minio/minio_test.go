package test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/bindings/s3"
	"github.com/balazsgrill/potatodrive/core"
	"github.com/balazsgrill/potatodrive/core/cfapi/filesystem"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Write func(p []byte) (n int, err error)

var _ io.Writer = Write(nil)

func (w Write) Write(p []byte) (n int, err error) {
	return w(p)
}

func stdlog(t *testing.T) Write {
	return func(p []byte) (n int, err error) {
		t.Log(string(p))
		return len(p), nil
	}
}

var MINIO_ACCESS_KEY = "minioadmin"
var MINIO_SECRET_KEY = "minioadmin"

func runAsync(t *testing.T, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = stdlog(t)
	cmd.Stderr = stdlog(t)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	return cmd
}

func runSync(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = stdlog(t)
	cmd.Stderr = stdlog(t)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

type testInstance struct {
	fsdir          string
	tempdir        string
	miniop         *exec.Cmd
	virtualization core.Virtualization
}

func setup(t *testing.T) *testInstance {
	instance := &testInstance{}
	instance.fsdir = t.TempDir()
	instance.tempdir = t.TempDir()
	minioworkdir := t.TempDir()
	instance.miniop = runAsync(t, "minio", "server", minioworkdir)
	time.Sleep(1 * time.Second)
	runSync(t, "mc", "alias", "set", "testminio", "http://localhost:9000", MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
	runSync(t, "mc", "mb", "testminio/test")
	return instance
}

func (instance *testInstance) start(t *testing.T) {
	config := s3.Config{
		Endpoint:  "localhost:9000",
		Region:    "us-east-1", // default region
		Bucket:    "test",
		KeyId:     MINIO_ACCESS_KEY,
		KeySecret: MINIO_SECRET_KEY,
		UseSSL:    false,
	}
	fs, err := config.ToFileSystem()
	if err != nil {
		t.Fatal(err)
	}
	instancecontext := bindings.InstanceContext{
		Logger: zerolog.New(zerolog.NewTestWriter(t)),
	}
	uid := uuid.NewMD5(uuid.UUID{}, []byte("test"))
	gid := core.BytesToGuid(uid[:])
	err = filesystem.RegisterRootPathSimple(*gid, instance.fsdir)
	if err != nil {
		t.Fatal(err)
	}
	instance.virtualization, err = filesystem.StartProjecting(instance.fsdir, fs, instancecontext.Logger)
	if err != nil {
		t.Fatal(err)
	}
}

func (instance *testInstance) Close() {
	if instance.virtualization != nil {
		instance.virtualization.Close()
	}
	filesystem.UnregisterRootPathSimple(instance.fsdir)
	instance.miniop.Process.Signal(os.Kill)
	instance.miniop.Wait()
}

func generateTestData(size int, seed string) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(seed[i%len(seed)])
	}
	return data
}

func TestDownloadingLargeFile(t *testing.T) {
	instance := setup(t)
	defer instance.Close()

	// generate data and upload it
	data := generateTestData(2*1024*1024, "2megabytesof2megabytes")
	inputfile := filepath.Join(instance.tempdir, "inputfile.dat")
	err := os.WriteFile(inputfile, data, 0644)
	if err != nil {
		t.Fatal(err)
	}
	runSync(t, "mc", "cp", inputfile, "testminio/test/root/inputfile.dat")

	instance.start(t)

	err = instance.virtualization.PerformSynchronization()
	if err != nil {
		t.Fatal(err)
	}

	outputfile := filepath.Join(instance.fsdir, "inputfile.dat")
	outputdata, err := os.ReadFile(outputfile)
	if err != nil {
		t.Fatal(err)
	}
	if string(outputdata) != string(data) {
		t.Fatal("data mismatch")
	}
}
