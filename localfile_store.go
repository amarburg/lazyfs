package lazyfs

import "os"
import "io"

import "fmt"
import "path/filepath"


type LocalFileStoreError struct {
	Err string
}

func (e LocalFileStoreError) Error() string {
	return e.Err
}

type LocalFileStore struct {
	source	FileSource
	root		string
	file    *os.File
}

func OpenLocalFileStore( source FileSource, root string ) (*LocalFileStore, error) {
	fs := LocalFileStore{ file: nil, root: root, source: source }
	return &fs, nil
}

// Load does the actual Lazy-loading of the file from the source to the
// local store.
func (fs *LocalFileStore) Load( )  error {
	if fs.file == nil {

		path := fs.root + fs.source.Path()

		os.MkdirAll( filepath.Dir(path), 0755 )

	f,err := os.Create( path )
	if err != nil {
		fmt.Println(err)
		return err
	}

	reader := fs.source.Reader()
	io.Copy( f, reader )

	f.Close();

fs.file,_ = os.Open(path)

	}
	return  nil
}

func (fs *LocalFileStore) ReadAt( p []byte, off int64) (n int, err error) {
	if err := fs.Load(); err != nil { return 0,err }

	return fs.file.ReadAt( p, off )
}

// func (fs *LocalFileStore) WriteAt(p []byte, off int64) (n int, err error) {
// 	return 0,nil
// }

func (fs *LocalFileStore) HasAt( p []byte, off int64 ) (n int, err error) {
	if err := fs.Load(); err != nil  {return 0,err}


	len := int64(cap( p ))
	sz,_ := fs.FileSize()

	switch {
		case (off + len) < sz: return int(len), nil
		case off > sz: return 0, LocalFileStoreError{"Offset beyond end of file"}
		case (off + len) > sz: return int(sz - off), nil
	}

	return 0, LocalFileStoreError{"Shouldn't get here"}
}

func (fs *LocalFileStore) FileSize() (int64,error) {
	stat,_ := fs.file.Stat()
	return stat.Size(),nil
}

func (fs *LocalFileStore) Reader() (io.Reader) {
	return fs.source.Reader()
}

func (fs *LocalFileStore) Path() string {
	return fs.source.Path()
}
