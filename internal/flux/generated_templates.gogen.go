// Code generated by vfsgen; DO NOT EDIT.

package flux

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// templates statically implements the virtual filesystem provided to vfsgen.
var templates = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"/flux-values.yaml.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "flux-values.yaml.tmpl",
			modTime:          time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			uncompressedSize: 736,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x52\xcb\x6e\xdb\x30\x10\xbc\xf3\x2b\x16\x39\x06\xa9\xf2\x40\x53\x14\xba\xb5\x3d\x18\x2d\x62\xa4\x70\xd0\x0f\x58\x53\x63\x89\x08\x1f\xea\x92\x04\xea\x1a\xf9\xf7\x42\xa4\x63\xa9\xf6\x4d\x3b\x3b\x2b\xce\xee\x8c\x20\x86\x2c\x1a\xb1\x55\x44\x82\xdf\x19\x31\x95\x6f\x22\x3d\xe6\x96\x1e\xef\x5c\x29\x1c\x5c\x90\x7d\x4b\x9f\x3e\xae\x8d\x52\x03\xac\x7b\x1e\x21\x9c\x82\x4c\x64\x2d\xe0\x84\x96\x76\x6c\x23\x94\xea\x4d\x9a\xd0\x2c\xb6\xa5\xc3\x81\x9a\x95\x49\xbf\x36\x4f\xf4\xf6\xa6\x88\xb6\xc2\x5e\x0f\x27\xfc\x6b\x29\x6b\xeb\x70\xf8\x40\x66\x57\xe0\x9f\x9c\x86\x58\xd1\x91\x53\xa5\xbf\x24\x31\xbe\x8f\x3f\x82\xf1\x0b\xce\xd5\xcd\xd5\x3c\x0d\xdf\xd5\x42\xc0\x5d\xf0\x76\x7f\x7a\x67\x03\xee\x9e\xbd\xdd\xd7\x76\x8e\x90\x59\x5a\x84\x54\xd8\xf2\x16\xb3\xe4\xa7\xa9\xaa\x0d\x6d\x5e\x5e\xcd\xd8\x52\x92\x0c\xa5\xde\x75\x7e\xd1\xb2\x41\x6f\x62\x92\xf2\x5b\x39\x7e\x4f\xab\xe3\x8f\xb6\xb9\xc3\x77\xc7\x3d\x5a\xe2\xd7\x28\x18\x43\x6c\xf8\x6f\x16\x68\x69\x4c\xb8\xbd\xbe\x71\x5a\x1a\x67\xb4\x84\x18\x76\xa9\xd1\xc1\xdd\x5e\x2b\x22\xd6\x52\xcf\x0f\xcf\x5b\x8b\xae\xca\x39\x7b\x6a\xb1\xab\x72\x70\x9a\xf5\x80\xae\x3a\xb8\xb0\x93\xc8\x1a\x67\xde\xed\x3c\x1a\x7a\x7f\x77\x74\x74\xe1\xe9\xc3\xe7\xb5\x29\xd8\xff\x01\x38\x8b\xc0\x3c\xf0\x78\xff\x30\xa5\x60\x94\xe0\x90\x06\xe4\xc2\x3f\xe9\xad\x47\x72\xec\xcd\x0e\x31\xad\xe0\xa7\x9c\x98\xe0\xeb\x26\xeb\x0b\xbc\x6c\x11\xf7\x5e\xaf\x58\xb6\xdc\xe3\x5b\xb0\x16\xba\x4c\xa8\xb3\x33\x5c\x30\xca\xac\xb6\x39\x26\xc8\x26\x58\x2c\xc3\x58\x74\xfc\x0b\x00\x00\xff\xff\xb0\x99\x13\xc0\xe0\x02\x00\x00"),
		},
		"/helm-operator-values.yaml.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "helm-operator-values.yaml.tmpl",
			modTime:          time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			uncompressedSize: 234,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\x8e\x39\x4b\x04\x41\x10\x85\xf3\xfa\x15\xc5\xc6\xae\x18\xa8\x41\xa7\x1a\x88\xa0\x81\x82\x79\xdb\xf3\x9c\x69\xe8\xcb\x3a\xc4\x65\xd9\xff\x2e\x33\xee\x66\x55\xef\xe3\x1d\x02\xed\x2e\x09\x1a\x88\x59\xf0\xed\x50\xdb\x6e\xe6\x34\x3c\xf0\xdd\x4d\xdd\x9e\x8a\xda\xe5\x10\xf8\xfe\xf6\x25\x13\xf9\x98\xa2\xe1\x61\x89\x62\x8f\x18\x1a\xd8\xc4\x41\xb4\xa0\xd4\xd5\xfb\x03\xd1\xdc\x9b\x06\x3e\x1e\xf9\xdd\x24\xb7\x59\x9f\x7b\x6e\x7c\xfd\x84\x52\x3f\xce\x94\x77\x57\x3b\x3e\x9d\x88\x86\xf4\x0a\x5b\xe0\x5b\x31\x5a\xfc\x2c\x98\x2e\x99\xa9\xb8\x1a\xe4\xad\x17\xac\x34\x09\xa2\xe1\x02\xe7\x6c\xab\xa8\xba\xfc\x4f\x56\x24\x81\xbd\xc6\x8a\xc0\x5f\xc5\x7f\xf7\x73\xb6\xfd\x84\x51\xfa\x81\xfe\x02\x00\x00\xff\xff\x63\x57\xe7\x06\xea\x00\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/flux-values.yaml.tmpl"].(os.FileInfo),
		fs["/helm-operator-values.yaml.tmpl"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
