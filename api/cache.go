package api

import (
	"crypto/md5"
	"encoding/hex"
)

func (a *Api) SaveCache() error {
	return a.cache.Save()
}

func (a *Api) RestoreCache() error {
	return a.cache.Restore()
}

func fileNameFromOrigin(origin string) string {
	sum := md5.Sum([]byte(origin))
	return hex.EncodeToString(sum[:])
}
