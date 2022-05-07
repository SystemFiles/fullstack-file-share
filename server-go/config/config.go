package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	CLOUD = "CLOUD"
	LOCAL = "LOCAL"
)	

type Config struct {
	Port int
	StorageMode string
	StorageDir string
	MaxFileCount int
	MaxFileSize int
	AutoClean bool
	CleanDays int
	CleanInterval int
	CORSAllowedOrigins []string
}

var once sync.Once
var instance Config

func Get() Config {
	once.Do(func() {
		p, err := strconv.Atoi(os.Getenv("FS_SERVER_PORT"))
		if err != nil {
			p = 5000
		}

		sm := strings.ToUpper(os.Getenv("FS_STORAGE_MODE"))
		if sm == "" {
			sm = LOCAL
		}

		sd := os.Getenv("FS_SERVER_STORAGE_DIR")
		if sd == "" {
			sd = "/Users/bensykes/Downloads"
		}

		mfc, err := strconv.Atoi(os.Getenv("FS_SERVER_MAX_FILE_COUNT"))
		if err != nil {
			mfc = 6
		}

		mfs, err := strconv.Atoi(os.Getenv("FS_SERVER_MAX_FILE_SIZE"))
		if err != nil {
			mfs = 1024
		}

		clean, err := strconv.ParseBool(os.Getenv("FS_SERVER_AUTO_CLEAN"))
		if err != nil {
			clean = true
		}

		cleanDays, err := strconv.Atoi(os.Getenv("FS_SERVER_CLEAN_DAYS"))
		if err != nil {
			cleanDays = 7
		}

		cleanInterval, err := strconv.Atoi(os.Getenv("FS_SERVER_CLEAN_INTERVAL"))
		if err != nil {
			cleanInterval = 7
		}

		allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
		if len(allowedOrigins) <= 1 || allowedOrigins == nil {
			allowedOrigins = []string{"*"}
		}

		instance = Config{
			Port: p,
			StorageMode: sm,
			StorageDir: sd,
			MaxFileCount: mfc,
			MaxFileSize: mfs,
			AutoClean: clean,
			CleanDays: cleanDays,
			CleanInterval: cleanInterval,
			CORSAllowedOrigins: allowedOrigins,
		}
	})

	return instance
}