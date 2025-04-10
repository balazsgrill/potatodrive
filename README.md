# PotatoDrive

## What?

Windows user-space locally-cached virtual file system client for windows 10+, capable of mounting remote folders as local. Can be used as convenient way to work with files stored in cloud storage or on a server, or as a way to backup and extend local storage. 

* Windows 10+
* Supports standard cloud or server storage backends without additional software
  * S3 (AWS, BackBlaze, Minio, etc..)
  * SFTP (SSH)
* Files are cached locally
* Multiple folder bindings on a single machine

## How?

Translates a [filesystem abstraction](https://github.com/spf13/afero) to a user space file system API provided by windows that can be either [Projected File System](https://learn.microsoft.com/en-us/windows/win32/projfs/projected-file-system) or [Cloud Files API](https://learn.microsoft.com/en-us/windows/win32/cfApi/cloud-files-api-portal).

## Why?

This software support two use cases:

* Extend the storage capacity of a local drive with a replicated cloud storage backend
  * Example 1: directly mount [AWS](https://aws.amazon.com/s3/) or [BackBlaze](https://www.backblaze.com/docs/cloud-storage-s3-compatible-api) buckets as local folders
  * Example 2: mount a self-hosted (optionally replicated) [Minio](https://min.io/) server
* Edit files stored on your SFTP server as if they were local files

## Installation

Installer is not (yet) available, latest binaries can be downloaded as [artifacts from successful builds](https://github.com/balazsgrill/potatodrive/actions).

If used, Projected File System service needs to be [enabled](https://learn.microsoft.com/en-us/windows/win32/projfs/enabling-windows-projected-file-system) (disabled by default):

```PowerShell
Enable-WindowsOptionalFeature -Online -FeatureName Client-ProjFS -NoRestart
```

## Configuration

Configuration is stored in Windows Registry, see [example.reg](example/potatodrive-minio.reg).

## Running

Once configured, just run the application. Logs are written to `%LOCALAPPDATA%\PotatoDrive`.

## Acknowledgements

This project could not have been possible without the following open source projects:
* [afero](https://github.com/spf13/afero)
* [walk](github.com/lxn/walk)
* [winrt-go](github.com/saltosystems/winrt-go)
