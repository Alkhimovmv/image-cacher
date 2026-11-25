package api

import (
	"testing"

	domainimage "image_cache/internal/domain/image"
)

type fakeImage struct {
	t *testing.T

	downloadFunc func(originURL string, fileName string) (domainimage.Image, error)
	readFunc     func(path string) ([]byte, error)
	removeFunc   func(path string) error
}

func (f *fakeImage) Download(originURL string, fileName string) (domainimage.Image, error) {
	f.t.Helper()
	if f.downloadFunc == nil {
		f.t.Fatalf("неожиданный вызов Image.Download(%q, %q)", originURL, fileName)
	}

	return f.downloadFunc(originURL, fileName)
}

func (f *fakeImage) Read(path string) ([]byte, error) {
	f.t.Helper()
	if f.readFunc == nil {
		f.t.Fatalf("неожиданный вызов Image.Read(%q)", path)
	}

	return f.readFunc(path)
}

func (f *fakeImage) Remove(path string) error {
	f.t.Helper()
	if f.removeFunc == nil {
		f.t.Fatalf("неожиданный вызов Image.Remove(%q)", path)
	}

	return f.removeFunc(path)
}

type fakeCache struct {
	t *testing.T

	getFunc     func(key string) (domainimage.Image, bool)
	setFunc     func(key string, value domainimage.Image, size uint) []domainimage.Image
	saveFunc    func() error
	restoreFunc func() error

	saveCalls    int
	restoreCalls int
}

func (f *fakeCache) Get(key string) (domainimage.Image, bool) {
	f.t.Helper()
	if f.getFunc == nil {
		f.t.Fatalf("неожиданный вызов Cache.Get(%q)", key)
	}

	return f.getFunc(key)
}

func (f *fakeCache) Set(key string, value domainimage.Image, size uint) []domainimage.Image {
	f.t.Helper()
	if f.setFunc == nil {
		f.t.Fatalf("неожиданный вызов Cache.Set(%q, %#v, %d)", key, value, size)
	}

	return f.setFunc(key, value, size)
}

func (f *fakeCache) Save() error {
	f.t.Helper()
	f.saveCalls++
	if f.saveFunc == nil {
		f.t.Fatal("неожиданный вызов Cache.Save()")
	}

	return f.saveFunc()
}

func (f *fakeCache) Restore() error {
	f.t.Helper()
	f.restoreCalls++
	if f.restoreFunc == nil {
		f.t.Fatal("неожиданный вызов Cache.Restore()")
	}

	return f.restoreFunc()
}
