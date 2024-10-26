package sftp

import (
	"errors"

	"github.com/balazsgrill/potatodrive/bindings/utils"
	sftpclient "github.com/pkg/sftp"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/afero/sftpfs"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	User       string `flag:"user,User name" reg:"User"`
	Password   string `flag:"password,Password" reg:"Password"`
	PrivateKey string `flag:"privatekey,PrivateKey" reg:"PrivateKey"`
	Host       string `flag:"host,Host:port" reg:"Host"`
	Basepath   string `flag:"basepath,Base path on remote server" reg:"Basepath"`
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return errors.New("host is mandatory")
	}
	if c.User == "" {
		return errors.New("user is mandatory")
	}
	if c.Password == "" && c.PrivateKey == "" {
		return errors.New("password or private key is mandatory")
	}
	return nil
}

func (c *Config) authFromKey() (ssh.AuthMethod, error) {
	key, err := ssh.ParseRawPrivateKey([]byte(c.PrivateKey))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

type configWithLogger struct {
	Config
	Logger zerolog.Logger
}

func (c *configWithLogger) Connect(onDisconnect func(error)) (afero.Fs, error) {
	var authmetods []ssh.AuthMethod
	if c.Password != "" {
		authmetods = append(authmetods, ssh.Password(c.Password))
	}
	if c.PrivateKey != "" {
		auth, err := c.authFromKey()
		if err != nil {
			c.Logger.Error().Err(err).Msg("SSH key auth failed")
		} else {
			authmetods = append(authmetods, auth)
		}
	}

	if len(authmetods) == 0 {
		return nil, errors.New("no valid authentication method is defined")
	}

	config := ssh.ClientConfig{
		User:            c.User,
		Auth:            authmetods,
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

func (c *Config) ToFileSystem(logger zerolog.Logger) (afero.Fs, error) {
	var remote afero.Fs
	cwithlogger := &configWithLogger{
		Config: *c,
		Logger: logger,
	}
	remote = &utils.ConnectingFs{
		Connect: cwithlogger.Connect,
	}
	if c.Basepath != "" {
		remote = utils.NewBasePathFs(remote, c.Basepath)
	}
	return remote, nil
}
