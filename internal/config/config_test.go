package config

import "testing"

func TestLoadReadsConfigFromYAMLFile(t *testing.T) {
	cfg, err := Load("../../config/config.yml")
	if err != nil {
		t.Fatalf("ожидали, что конфиг прочитается без ошибки, получили %v", err)
	}

	if cfg.HTTP.Addr != ":3001" {
		t.Fatalf("ожидали http.addr = :3001, получили %s", cfg.HTTP.Addr)
	}

	if cfg.Img.Folder != "app-data/images" {
		t.Fatalf("ожидали img.folder = app-data/images, получили %s", cfg.Img.Folder)
	}

	if cfg.Img.Timeout != 5 {
		t.Fatalf("ожидали img.timeout = 5, получили %d", cfg.Img.Timeout)
	}

	if cfg.Cache.MaxSize != 52428800 {
		t.Fatalf("ожидали cache.max_size = 52428800, получили %d", cfg.Cache.MaxSize)
	}

	if cfg.Cache.Path != "app-data/cache/data.json" {
		t.Fatalf("ожидали cache.path = app-data/cache/data.json, получили %s", cfg.Cache.Path)
	}
}
