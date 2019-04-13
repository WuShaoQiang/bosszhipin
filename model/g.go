package model

import (
	"path/filepath"

	"github.com/WuShaoQiang/crawler/boss/config"
	"github.com/jinzhu/gorm"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	db       *gorm.DB
	logger   = log.New()
	basePath = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss"
)

func init() {
	setLogger()
}

func setLogger() {
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: filepath.Join(basePath + "/log/debug.log"),
		logrus.InfoLevel:  filepath.Join(basePath + "/log/info.log"),
		logrus.WarnLevel:  filepath.Join(basePath + "/log/warn.log"),
	}
	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}

// SetDB func
func SetDB(database *gorm.DB) {
	db = database
}

// ConnectToDB func
func ConnectToDB() *gorm.DB {
	connectingStr := config.GetMysqlConnectingString()
	log.Println("Connet to db...")
	db, err := gorm.Open("mysql", connectingStr)
	if err != nil {
		log.Fatalf("Failed to connect database : %s", err)
	}
	db.SingularTable(true)
	return db
}
