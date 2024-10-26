typedef i32 FileHandle
typedef i32 FileMode
typedef i64 Timestamp

exception FilesystemException {
    1: string message
    2: bool isnotexists
    3: bool eof
  }

struct FileInfo {
  1: string fname
  2: i64 fsize
  3: FileMode fmode
  4: Timestamp ftime
  5: bool fisDir
}

service Filesystem{
    FileHandle create(1:string name) throws(1:FilesystemException error)
    void mkdir(1:string path, 2:FileMode perm) throws(1:FilesystemException error)
    void mkdirAll(1:string path, 2:FileMode perm) throws(1:FilesystemException error)
    FileHandle open(1:string name) throws(1:FilesystemException error)
    FileHandle openFile(1:string name, 2:i32 flag, 3:FileMode perm) throws(1:FilesystemException error)
    void remove(1:string name) throws(1:FilesystemException error)
    void removeAll(1:string name) throws(1:FilesystemException error)
    void rename(1:string oldname, 2:string newname) throws(1:FilesystemException error)
    FileInfo stat(1:string name) throws(1:FilesystemException error)
    string name()
    void chmod(1:string name, 2:FileMode mode) throws(1:FilesystemException error)
    void chown(1:string name, 2:i32 uid, 3:i32 gid) throws(1:FilesystemException error)
    void chtimes(1:string name, 2:Timestamp atime, 3:Timestamp mtime) throws(1:FilesystemException error)

    // File operations
    // Closer
    void fclose(1:FileHandle file) throws(1:FilesystemException error)
    // Reader
    binary fread(1:FileHandle file, 2:i64 bufferSize) throws(1:FilesystemException error)
    // ReaderAt
    binary freadAt(1:FileHandle file, 2:i64 bufferSize, 3:i64 offset) throws(1:FilesystemException error)
    // Seeker
    i64 fseek(1:FileHandle file, 2:i64 offset, 3:i32 whence) throws(1:FilesystemException error)
    // Writer
    i32 fwrite(1:FileHandle file, 2:binary buffer) throws(1:FilesystemException error)
    // WriterAt
    i32 fwriteAt(1:FileHandle file, 2:binary buffer, 3:i64 offset) throws(1:FilesystemException error)

    string fname(1:FileHandle file)
    list<FileInfo> freaddir(1:FileHandle file, 2:i32 count) throws(1:FilesystemException error)
    list<string> freaddirnames(1:FileHandle file, 2:i32 count) throws(1:FilesystemException error)
    FileInfo fstat(1:FileHandle file) throws(1:FilesystemException error)
    void fsync(1:FileHandle file) throws(1:FilesystemException error)
    void ftruncate(1:FileHandle file, 2:i64 size) throws(1:FilesystemException error)
    i32 fwriteString(1:FileHandle file, 2:string value) throws(1:FilesystemException error)
}