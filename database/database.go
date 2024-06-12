package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lulzshadowwalker/pupsik/config"
	md "github.com/lulzshadowwalker/pupsik/database/.gen/postgres/public/model"
	. "github.com/lulzshadowwalker/pupsik/database/.gen/postgres/public/table"
	"github.com/lulzshadowwalker/pupsik/database/model"
	"github.com/lulzshadowwalker/pupsik/types"
)

var DB *sql.DB

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

	DB = database
}

func GetAccountByUserID(ctx context.Context, id uuid.UUID) (types.Account, error) {
	var dest md.Account
	stmt := SELECT(Account.AllColumns).
		FROM(Account).
		WHERE(Account.UserID.EQ(UUID(id)))

	if err := stmt.QueryContext(ctx, DB, &dest); err != nil {
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

	var db qrm.Executable = DB
	if tx != nil {
		db = tx
	}
	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return types.Account{}, fmt.Errorf("failed to update account because %w", err)
	}

	return account, nil
}

func GetImagesByUserID(ctx context.Context, id uuid.UUID, tx *sql.Tx) ([]types.Image, error) {
	var db qrm.Queryable = DB
	if tx != nil {
		db = tx
	}

	var dest []model.DBImage
	if err := Image.SELECT(Image.AllColumns).
		FROM(Image).
		WHERE(Image.UserID.EQ(postgres.UUID(id))).QueryContext(ctx, db, &dest); err != nil {
		return nil, fmt.Errorf("failed to fetch images by user id because %w", err)
	}

	images := make([]types.Image, len(dest))
	for i := range dest {
		images[i] = dest[i].ToEntity()
	}

	return images, nil
}

func GetImageByID(ctx context.Context, id int, tx *sql.Tx) (types.Image, error) {
	var db qrm.Queryable = DB
	if tx != nil {
		db = tx
	}

	var dest model.DBImage
	if err := Image.SELECT(Image.AllColumns).
		FROM(Image).
		WHERE(Image.ID.EQ(postgres.Int64(int64(id)))).QueryContext(ctx, db, &dest); err != nil {
		return types.Image{}, fmt.Errorf("failed to fetch image by id because %w", err)
	}

	return dest.ToEntity(), nil
}

func StoreImage(ctx context.Context, img types.Image, tx *sql.Tx) (types.Image, error) {
	var db qrm.Queryable = DB
	if tx != nil {
		db = tx
	}

	var dest model.DBImage
	if err := Image.INSERT(
		Image.UserID,
		Image.BatchID,
		Image.Prompt,
	).VALUES(
		img.UserID,
		img.BatchID,
		img.Prompt,
	).RETURNING(Image.ID, Image.Status).
		QueryContext(ctx, db, &dest); err != nil {
		return types.Image{}, fmt.Errorf("failed to store image because %w", err)
	}

	img.Status = types.ImageStatus(dest.Image.Status)
	img.ID = int(dest.Image.ID)
	return img, nil
}

func GetImagesByBatchID(ctx context.Context, id uuid.UUID, tx *sql.Tx) ([]types.Image, error) {
	var db qrm.Queryable = DB
	if tx != nil {
		db = tx
	}

	var dest []model.DBImage
	if err := Image.SELECT(Image.AllColumns).
		FROM(Image).
		WHERE(Image.BatchID.EQ(postgres.UUID(id))).QueryContext(ctx, db, &dest); err != nil {
		return nil, fmt.Errorf("failed to fetch images by batch id because %w", err)
	}

	images := make([]types.Image, len(dest))
	for i := range dest {
		images[i] = dest[i].ToEntity()
	}

	return images, nil
}

func UpdateImage(ctx context.Context, img types.Image, tx *sql.Tx) (types.Image, error) {
	var db qrm.Queryable = DB
	if tx != nil {
		db = tx
	}

	var dest model.DBImage
	if err := Image.UPDATE(
		Image.URL,
		Image.Status,
	).
		SET(
			img.URL,
			img.Status,
		).
		WHERE(
			Image.ID.EQ(postgres.Int64(int64(img.ID))),
		).RETURNING(Image.AllColumns).
		QueryContext(ctx, db, &dest); err != nil {
		return types.Image{}, fmt.Errorf("failed to update image because %w", err)
	}

	return dest.ToEntity(), nil
}
