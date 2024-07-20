package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/balazsgrill/projfero/filesystem"
	s3 "github.com/fclairamb/afero-s3"
)

func main() {
	endpoint := flag.String("endpoint", "", "S3 endpoint")
	accessKeyID := flag.String("keyid", "", "Access Key ID")
	secretAccessKey := flag.String("secred", "", "Access Key Secret")
	useSSL := flag.Bool("useSSL", false, "Use SSL encryption for S3 connection")
	region := flag.String("region", "", "Region")
	bucket := flag.String("bucket", "", "Bucket")
	localpath := flag.String("localpath", "", "Local folder")
	flag.Parse()

	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(*region),
		Endpoint:         aws.String(*endpoint),
		DisableSSL:       aws.Bool(!*useSSL),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(*accessKeyID, *secretAccessKey, ""),
	})

	// Initialize the file system
	fs := s3.NewFs(*bucket, sess)
	fs.MkdirAll("root", 0777)
	rootfs := NewBasePathFs(fs, "root")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	closer, err := filesystem.StartProjecting(*localpath, rootfs)
	if err != nil {
		log.Panic(err)
	}

	t := time.NewTicker(30 * time.Second)
	go func() {
		for range t.C {
			err = closer.PerformSynchronization()
			if err != nil {
				log.Panic(err)
			}
		}
	}()

	<-c
	t.Stop()
	closer.Close()
	os.Exit(1)
}
