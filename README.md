# PotatoDrive

Windows user-space locally-cached virtual file system client. Based on [Projected File System](https://learn.microsoft.com/en-us/windows/win32/projfs/projected-file-system) service available on Windows 10+ and works similar to OneDrive, using the local directory as a cache of the remote state.

:construction: PotatoDrive is in early development phase, don't trust it with data you're not willing to lose. Keep your backups safe (as always)!

Features:

* Windows 10+
* Supports standard S3 (AWS, BackBlaze, Minio, etc..) as backend without additional server
* Files are cached locally
* Multiple folder bindings on a single machine

## Installation

Installer is not (yet) available, latest binaries can be downloaded as [artifacts from successful builds](https://github.com/balazsgrill/potatodrive/actions).

To work, Projected File System service needs to be [enabled](https://learn.microsoft.com/en-us/windows/win32/projfs/enabling-windows-projected-file-system) (disabled by default):

```PowerShell
Enable-WindowsOptionalFeature -Online -FeatureName Client-ProjFS -NoRestart
```

## Configuration

Configuration is stored in Windows Registry, see [example.reg](example/potatodrive-minio.reg).

## Running

Once configured, 

# Development

* [afero](https://github.com/spf13/afero) is used as backend interface, meaning that any 

Windows user-space virtual file system binding to [afero](https://github.com/spf13/afero) file systems backend. Makes it possible to implement virtual file systems for windows backed by any kind of user-space implementation of virtual file system.
