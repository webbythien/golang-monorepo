package must

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/webbythien/monorepo/pkg/l"
	"github.com/webbythien/monorepo/pkg/watcher"
	"github.com/webbythien/monorepo/sdk/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxDBIdleConns  = 10
	maxDBOpenConns  = 100
	maxConnLifeTime = 30 * time.Minute
)

func ConnectPostgreSQL(cfg *conf.PostgreSQL) *gorm.DB {
	if cfg.RDS != nil && cfg.RDS.IsEnabled {
		ll.Info("Using AWS context for PostgreSQL connection")
		sess := session.Must(session.NewSession())
		creds := sess.Config.Credentials

		// awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
		// if err != nil {
		// 	panic("configuration error: " + err.Error())
		// }
		dbEndpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		// authenticationToken, err := auth.BuildAuthToken(
		// 	context.TODO(), dbEndpoint, cfg.RDS.Region, cfg.Username, aws.CredentialsProvider(creds))
		// if err != nil {
		// 	panic("failed to create authentication token: " + err.Error())
		// }
		authenticationToken, err := rdsutils.BuildAuthToken(dbEndpoint, cfg.RDS.Region, cfg.Username, creds)
		if err != nil {
			panic(err)
		}

		cfg.Password = authenticationToken
	}
	logMode := logger.Error
	if os.Getenv("POSTGRESQL_DEBUG") == "true" {
		logMode = logger.Info
	}
	db, err := gorm.Open(postgres.Open(cfg.FormatDSN()), &gorm.Config{
		//Logger: zapgorm.New(ll.Logger).LogMode(logMode),
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		ll.Fatal("Error open sql db", l.Error(err))
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(cfg.Database))); err != nil {
		ll.Fatal("Error use otelgorm", l.Error(err))
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		ll.Fatal("Error querying SELECT 1", l.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		ll.Fatal("Error get sql DB", l.Error(err))
	}
	watcher.RegisterCleanFunc(func() {
		ll.Info("Closing postgresql db connection")
		err = sqlDB.Close()
		if err != nil {
			ll.Info("Failed to close postgresql db", l.Error(err))
		}
	})
	sqlDB.SetMaxIdleConns(maxDBIdleConns)
	sqlDB.SetMaxOpenConns(maxDBOpenConns)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	return db
}
