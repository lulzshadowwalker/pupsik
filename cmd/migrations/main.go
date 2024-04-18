package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/lulzshadowwalker/pupsik/config"
	"github.com/serenize/snaker"
)

func main() {
	port, err := strconv.Atoi(config.GetDatabasePort())
	if err != nil {
		log.Fatalf("failed to read database port because %s", err)
	}

	dbConnection := postgres.DBConnection{
		Host:       config.GetDatabaseHost(),
		Port:       port,
		User:       config.GetDatabaseUser(),
		Password:   config.GetDatabasePassword(),
		DBName:     config.GetDatabaseName(),
		SchemaName: "public",
		SslMode:    "require",
	}

	err = postgres.Generate(
		"./database/.gen/",
		dbConnection,
		template.Default(postgres2.Dialect).
			UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
				return template.DefaultSchema(schemaMetaData).
					UseModel(template.DefaultModel().
						UseTable(func(table metadata.Table) template.TableModel {
							return template.DefaultTableModel(table).
								UseField(func(columnMetaData metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(columnMetaData)
									return defaultTableModelField.UseTags(
										fmt.Sprintf(`json:"%s,omitempty"`, snaker.SnakeToCamelLower(columnMetaData.Name)),
									)
								})
						}),
					)
			}),
	)
	if err != nil {
		log.Fatalf("failed to migrate data because %s", err)
	}
}
