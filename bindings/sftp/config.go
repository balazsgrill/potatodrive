package sftp

import (
	"errors"

	"github.com/balazsgrill/potatodrive/bindings/utils"
	sftpclient "github.com/pkg/sftp"
	"github.com/spf13/afero"
	"github.com/spf13/afero/sftpfs"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	User     string `flag:"user,User name" reg:"User"`
	Password string `flag:"password,Password" reg:"Password"`
	Host     string `flag:"host,Host:port" reg:"Host"`
	Basepath string `flag:"basepath,Base path on remote server" reg:"Basepath"`
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is mandatory")
	}
	if c.User == "" {
		return errors.New("user is mandatory")
	}
	if c.Password == "" {
		return errors.New("password is mandatory")
	}
	return nil
}

func (c *Config) Connect(onDisconnect func(error)) (afero.Fs, error) {
	config := ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", c.Host, &config)
	if err != nil {
		return nil, err
	}

	client, err := sftpclient.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	go func() {
		err := conn.Wait()
		client.Close()
		onDisconnect(err)
	}()

	return sftpfs.New(client), nil
}

func (c *Config) ToFileSystem() (afero.Fs, error) {
	var remote afero.Fs
	remote = &utils.ConnectingFs{
		Connect: c.Connect,
	}
	if c.Basepath != "" {
		remote = utils.NewBasePathFs(remote, c.Basepath)
	}
	return remote, nil
}
