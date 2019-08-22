package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/brunoksato/golang-boilerplate/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	elastic "github.com/olivere/elastic"
	"go.elastic.co/apm/module/apmgorm"
)

var CONFIGURATIONS map[string]string = map[string]string{
	"SERVER_ENV":            "development",
	"SERVER_NAME":           "SERVER",
	"DATABASE_URL":          "dbname=server sslmode=disable",
	"AWS_ACCESS_KEY":        "",
	"AWS_ACCESS_KEY_SECRET": "",
	"AWS_REGION":            "us-east-1",
	"AWS_BUCKET":            "development.company.asset",
	"ES_HOST":               "localhost:9200",
	"LOGGER_LEVEL":          "info",
	"SENDGRID_KEY":          "key_sendgrid",
	"SENDGRID_USER":         "support@company.com",
	"JWT_KEY_SIGNIN":        "you_secret_key",
	"JWT_KEY_EMAIL":         "you_secret_key_email",
	"JWT_TOKEN_EXPIRATION":  "72",
}

func Init() {
	if os.Getenv("SERVER_ENV") == "" {
		os.Setenv("SERVER_ENV", CONFIGURATIONS["SERVER_ENV"])
	}
	if os.Getenv("SERVER_NAME") == "" {
		os.Setenv("SERVER_NAME", CONFIGURATIONS["SERVER_NAME"])
	}
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", CONFIGURATIONS["DATABASE_URL"])
	}
	if os.Getenv("AWS_ACCESS_KEY") == "" {
		os.Setenv("AWS_ACCESS_KEY", CONFIGURATIONS["AWS_ACCESS_KEY"])
	}
	if os.Getenv("AWS_ACCESS_KEY_SECRET") == "" {
		os.Setenv("AWS_ACCESS_KEY_SECRET", CONFIGURATIONS["AWS_ACCESS_KEY_SECRET"])
	}
	if os.Getenv("AWS_REGION") == "" {
		os.Setenv("AWS_REGION", CONFIGURATIONS["AWS_REGION"])
	}
	if os.Getenv("AWS_BUCKET") == "" {
		os.Setenv("AWS_BUCKET", CONFIGURATIONS["AWS_BUCKET"])
	}
	if os.Getenv("ES_HOST") == "" {
		os.Setenv("ES_HOST", CONFIGURATIONS["ES_HOST"])
	}
	if os.Getenv("LOGGER_LEVEL") == "" {
		os.Setenv("LOGGER_LEVEL", CONFIGURATIONS["LOGGER_LEVEL"])
	}
	if os.Getenv("SENDGRID_KEY") == "" {
		os.Setenv("SENDGRID_KEY", CONFIGURATIONS["SENDGRID_KEY"])
	}
	if os.Getenv("SENDGRID_USER") == "" {
		os.Setenv("SENDGRID_USER", CONFIGURATIONS["SENDGRID_USER"])
	}
	if os.Getenv("JWT_KEY_SIGNIN") == "" {
		os.Setenv("JWT_KEY_SIGNIN", CONFIGURATIONS["JWT_KEY_SIGNIN"])
	}
	if os.Getenv("JWT_KEY_EMAIL") == "" {
		os.Setenv("JWT_KEY_EMAIL", CONFIGURATIONS["JWT_KEY_EMAIL"])
	}
	if os.Getenv("JWT_TOKEN_EXPIRATION") == "" {
		os.Setenv("JWT_TOKEN_EXPIRATION", CONFIGURATIONS["JWT_TOKEN_EXPIRATION"])
	}
}

func InitDB() *gorm.DB {
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", CONFIGURATIONS["DATABASE_URL"])
	}

	db, err := apmgorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fmt.Sprintf("Initialized read-write database connection pool: %s", os.Getenv("DATABASE_URL")))
	return db
}

func InitElasticSearchAndLogger() (client *elastic.Client) {
	logrus.SetFormatter(&log.LogstashFormatter{})
	// Acceptable values are:
	// debug, info, warn, error, fatal, panic
	level := os.Getenv("LOGGER_LEVEL")
	if level != "" {
		l, err := logrus.ParseLevel(level)
		if err == nil {
			logrus.SetLevel(l)
		} else {
			fmt.Println("Error with log level configuraion:", err)
		}
	}

	appName := fmt.Sprintf("%s-%s", os.Getenv("NAME"), os.Getenv("ENV"))
	esHostname := os.Getenv("ES_HOST")
	thisHostname, _ := os.Hostname()

	if esHostname != "" {
		var esURL string
		if strings.Contains(esHostname, "127.0.0.1") {
			esURL = fmt.Sprintf("http://%s", esHostname)
		} else {
			esURL = fmt.Sprintf("https://%s", esHostname)
		}
		fmt.Println(fmt.Sprintf("Configuring elasticsearch logging: %s", esURL))
		client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esURL))
		if err != nil {
			fmt.Println(fmt.Sprintf("Error configuring elasticsearch logging: %s", err.Error()))
		} else {
			now := time.Now()
			appName = fmt.Sprintf("logstash-%d.%d.%d", now.Year(), now.Month(), now.Day())
			hook, err := log.NewElasticHook(client, thisHostname, logrus.DebugLevel, appName)
			if err == nil {
				logrus.AddHook(hook)
			} else {
				fmt.Println(fmt.Sprintf("Error configuring logger for elastic search: %s", err.Error()))
			}
		}

		return client
	}

	return
}
