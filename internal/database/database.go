package database

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"pairs-trading-backend/internal/models"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPass, dbName)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseConfig: %v", err)
	}

	config.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}

	dbPool := stdlib.OpenDB(*config)

	dbPool.SetMaxOpenConns(25)
	dbPool.SetMaxIdleConns(5)
	dbPool.SetConnMaxLifetime(5 * time.Minute)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbPool,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %v", err)
	}

	err = dbPool.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to Cloud SQL database")

	err = gormDB.Exec("SET search_path TO gold, public").Error
	if err != nil {
		return nil, fmt.Errorf("failed to set search path: %v", err)
	}

	migrator := gormDB.Migrator()
	tables := []interface{}{
		&models.ETFDailyOHLC{},
		&models.SuggestedPair{},
		&models.NewsMention{},
		&models.Ticker{},
	}

	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Printf("Table %T does not exist", table)
		}
	}

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate User model: %v", err)
	}

	return gormDB, nil
}
