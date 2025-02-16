package asset

import (
	"crypto/sha256"
	"fmt"
)

type Asset struct {
	Name    string
	content []byte `json:"-"`
	Hashsum [sha256.Size]byte
}

func NewAsset(name string) *Asset {
	return &Asset{
		Name: name,
	}
}

func (a *Asset) Equal(other *Asset) bool {
	return a.Hashsum == other.Hashsum
}

func (a *Asset) SetContent(content []byte) {
	a.content = content
	a.Hashsum = sha256.Sum256(a.content)
}

func (a *Asset) String() string {
	return fmt.Sprintf("Name=%s", a.Name)
}
