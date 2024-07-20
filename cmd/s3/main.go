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
	"golang.org/x/sys/windows/registry"
)

func main() {
	regkey := flag.String("regkey", "", "Registry key that holds configuration")
	endpoint := flag.String("endpoint", "", "S3 endpoint")
	accessKeyID := flag.String("keyid", "", "Access Key ID")
	secretAccessKey := flag.String("secred", "", "Access Key Secret")
	useSSL := flag.Bool("useSSL", false, "Use SSL encryption for S3 connection")
	region := flag.String("region", "", "Region")
	bucket := flag.String("bucket", "", "Bucket")
	localpath := flag.String("localpath", "", "Local folder")
	flag.Parse()

	if *regkey != "" {
		key, err := registry.OpenKey(registry.CURRENT_USER, *regkey, registry.QUERY_VALUE)
		if err != nil {
			log.Panic(err)
		}
		*endpoint, _, err = key.GetStringValue("Endpoint")
		if err != nil {
			log.Panic(err)
		}
		*accessKeyID, _, err = key.GetStringValue("KeyID")
		if err != nil {
			log.Panic(err)
		}
		*secretAccessKey, _, err = key.GetStringValue("KeySecret")
		if err != nil {
			log.Panic(err)
		}
		useSSLint, _, err := key.GetIntegerValue("UseSSL")
		if err != nil {
			log.Panic(err)
		}
		*useSSL = useSSLint != 0
		*region, _, err = key.GetStringValue("Region")
		if err != nil {
			log.Panic(err)
		}
		*bucket, _, err = key.GetStringValue("Bucket")
		if err != nil {
			log.Panic(err)
		}
		*localpath, _, err = key.GetStringValue("Directory")
		if err != nil {
			log.Panic(err)
		}
	}

	if *endpoint == "" || *accessKeyID == "" || *secretAccessKey == "" || *region == "" || *bucket == "" || *localpath == "" {
		log.Panic("Missing configuration")
	}

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
