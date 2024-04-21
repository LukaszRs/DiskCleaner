package models

type FileDef struct {
	ID         string
	IsDir      bool
	Path       string
	Filename   string
	Size       int64
	CreatedAt  string
	ModifiedAt string
	Checksum   string
	Hash       string
}
