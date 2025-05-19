package evaluator

import (
	"fmt"
	"hash/fnv"
)

type FileValue struct {
	FileID      string
	FilePath    string
	FileName    string
	FileSize    int64
	ContentType string
}

func NewFileValue(fileID, filePath, fileName, contentType string, fileSize int64) *FileValue {
	return &FileValue{
		FileID:      fileID,
		FilePath:    filePath,
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
	}
}

func (s *FileValue) Debug() string {
	return fmt.Sprintf("<file: %s, size: %d, contentType: %s>", s.FileName, s.FileSize, s.ContentType)
}

func (s *FileValue) Type() ObjectType {
	return FileObject
}

func (s *FileValue) HashKey() HashKey {

	h := fnv.New64a()
	_, _ = h.Write([]byte(s.FileID))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
