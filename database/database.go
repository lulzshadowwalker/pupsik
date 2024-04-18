package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lulzshadowwalker/pupsik/config"
	md "github.com/lulzshadowwalker/pupsik/database/.gen/postgres/public/model"
	. "github.com/lulzshadowwalker/pupsik/database/.gen/postgres/public/table"
	"github.com/lulzshadowwalker/pupsik/types"
)

var db *sql.DB

func init() {
	port, err := strconv.Atoi(config.GetDatabasePort())
	if err != nil {
		log.Fatalf("failed to read database port because %s", err)
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		config.GetDatabaseHost(),
		port,
		config.GetDatabaseUser(),
		config.GetDatabaseName(),
		"require",
		config.GetDatabasePassword(),
	)

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Errorf("failed to open connect to database because %w", err))
	}

	err = database.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping the database because %w", err))
	}

	db = database
}

func GetAccountByUserID(ctx context.Context, id uuid.UUID) (types.Account, error) {
	var dest md.Account
	stmt := SELECT(Account.AllColumns).
		FROM(Account).
		WHERE(Account.UserID.EQ(UUID(id)))

	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		return types.Account{}, fmt.Errorf("failed to get account by id because %w", err)
	}

	return types.Account{
		Username: dest.UserName,
		ID:       int(dest.ID),
		UserID:   dest.UserID,
	}, nil
	// return dest.ToEntity(), nil
}

func CreateAccount(ctx context.Context, account types.Account, tx *sql.Tx) (types.Account, error) {
	stmt := Account.INSERT(
		Account.UserName,
		Account.UserID,
	).VALUES(
		account.Username,
		account.UserID,
	)

	var db qrm.Executable = db
	if tx != nil {
		db = tx
	}
	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return types.Account{}, fmt.Errorf("failed to update account because %w", err)
	}

	return account, nil
}
