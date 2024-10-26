package server

import (
	"context"
	"os"

	"github.com/balazsgrill/potatodrive/bindings/proxy"
)

func (fs *FilesystemServer) Fclose(ctx context.Context, file proxy.FileHandle) (_err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return os.ErrInvalid
	}
	err := f.Close()
	if err != nil {
		return ewrap(err)
	}
	delete(fs.openfiles, file)
	return nil
}

func (fs *FilesystemServer) Fname(ctx context.Context, file proxy.FileHandle) (_r string, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return "", os.ErrInvalid
	}
	return f.Name(), nil
}

func (fs *FilesystemServer) Fread(ctx context.Context, file proxy.FileHandle, bufferSize int64) (_r []byte, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return nil, os.ErrInvalid
	}
	buffer := make([]byte, bufferSize)
	n, err := f.Read(buffer)
	return buffer[:n], ewrap(err)
}

func (fs *FilesystemServer) FreadAt(ctx context.Context, file proxy.FileHandle, bufferSize int64, offset int64) (_r []byte, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return nil, os.ErrInvalid
	}
	buffer := make([]byte, bufferSize)
	n, err := f.ReadAt(buffer, offset)
	return buffer[:n], ewrap(err)
}

// Freaddir implements proxy.Filesystem.
func (fs *FilesystemServer) Freaddir(ctx context.Context, file proxy.FileHandle, count int32) (_r []*proxy.FileInfo, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return nil, os.ErrInvalid
	}
	infos, err := f.Readdir(int(count))

	wrappedInfos := make([]*proxy.FileInfo, len(infos))
	for i, info := range infos {
		wrappedInfos[i] = wrapFileInfo(info)
	}

	return wrappedInfos, ewrap(err)
}

func (fs *FilesystemServer) Freaddirnames(ctx context.Context, file proxy.FileHandle, count int32) (_r []string, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return nil, os.ErrInvalid
	}
	names, err := f.Readdirnames(int(count))
	return names, ewrap(err)
}

func (fs *FilesystemServer) Fseek(ctx context.Context, file proxy.FileHandle, offset int64, whence int32) (_r int64, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return 0, os.ErrInvalid
	}
	pos, err := f.Seek(offset, int(whence))
	return pos, ewrap(err)
}

func (fs *FilesystemServer) Fstat(ctx context.Context, file proxy.FileHandle) (_r *proxy.FileInfo, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return nil, os.ErrInvalid
	}
	info, err := f.Stat()
	if err != nil {
		return nil, ewrap(err)
	}
	return wrapFileInfo(info), nil
}

func (fs *FilesystemServer) Fsync(ctx context.Context, file proxy.FileHandle) (_err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return os.ErrInvalid
	}
	return ewrap(f.Sync())
}

func (fs *FilesystemServer) Ftruncate(ctx context.Context, file proxy.FileHandle, size int64) (_err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return os.ErrInvalid
	}
	return ewrap(f.Truncate(size))
}

func (fs *FilesystemServer) Fwrite(ctx context.Context, file proxy.FileHandle, buffer []byte) (_r int32, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return 0, os.ErrInvalid
	}
	r, err := f.Write(buffer)
	return int32(r), ewrap(err)
}

func (fs *FilesystemServer) FwriteAt(ctx context.Context, file proxy.FileHandle, buffer []byte, offset int64) (_r int32, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return 0, os.ErrInvalid
	}
	r, err := f.WriteAt(buffer, offset)
	return int32(r), ewrap(err)
}

func (fs *FilesystemServer) FwriteString(ctx context.Context, file proxy.FileHandle, value string) (_r int32, _err error) {
	f, opened := fs.openfiles[file]
	if !opened {
		return 0, os.ErrInvalid
	}
	r, err := f.WriteString(value)
	return int32(r), ewrap(err)
}
