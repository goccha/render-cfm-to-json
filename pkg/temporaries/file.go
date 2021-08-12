package temporaries

import (
	"github.com/goccha/render-cfm-to-json/pkg/debug"
	"golang.org/x/xerrors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	name string
	file *os.File
}

func (f *File) Name() string {
	return f.name
}
func (f *File) Write(b []byte) (n int, err error) {
	if n, err = f.file.Write(b); err != nil {
		err = xerrors.Errorf(": %w", err)
	}
	return
}

func Open(filename string) (*File, error) {
	if tmpFile, filePath, err := createTempFile(filename); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else {
		return &File{
			name: filePath,
			file: tmpFile,
		}, nil
	}
}

func fileName(name string) string {
	_, filename := filepath.Split(name)
	if index := strings.LastIndex(filename, "."); index > 0 {
		filename = filename[0:index]
	}
	filename = filename + strconv.Itoa(rand.Int()) + ".json"
	return filename
}

func createTempFile(name string) (tmpFile *os.File, filePath string, err error) {
	filename := fileName(name)
	debug.Print("filename", filename)
	v := os.Getenv("GITHUB_WORKSPACE")
	work := filepath.Join(v, filename)
	if tmpFile, err = os.Create(work); err != nil {
		return nil, "", xerrors.Errorf(": %w", err)
	}
	debug.Print("tmpFile", tmpFile.Name())
	workDir := os.Getenv("RUNNER_WORKSPACE")
	repo := os.Getenv("GITHUB_REPOSITORY")
	repo = strings.Split(repo, "/")[1]
	dir := filepath.Join(workDir, repo)
	filePath = filepath.Join(dir, filename)
	return
}
