package mumbler

import (
	"github.com/layeh/gumble/gumbleffmpeg"
	"io"
	"io/ioutil"
)

type Source interface {
	GetSource() gumbleffmpeg.Source
}

type fileSource struct {
	source string
}

func NewFileSource(path string) *fileSource {
	return &fileSource{path}
}

func (source *fileSource) GetSource() gumbleffmpeg.Source {
	return gumbleffmpeg.SourceFile(source.source)
}

type readerSource struct {
	source io.ReadCloser
}

func NewReaderSource(reader io.Reader) *readerSource {
	return &readerSource{ioutil.NopCloser(reader)}
}

func NewReadCloserSource(reader io.ReadCloser) *readerSource {
	return &readerSource{reader}
}

func (source *readerSource) GetSource() gumbleffmpeg.Source {
	return gumbleffmpeg.SourceReader(source.source)
}
