package henge

import (
	"fmt"
	"strings"

	"github.com/btnguyen2k/prom"
)

// CreateTable generates and executes "CREATE TABLE" statement
func CreateTable(sqlc *prom.SqlConnect, tableName string, ifNotExist bool, colDef map[string]string, colNames, pk []string) error {
	template := "CREATE TABLE %s %s (%s%s)"
	partIfNotExists := ""
	if ifNotExist {
		partIfNotExists = "IF NOT EXISTS"
	}
	partColDef := make([]string, 0)
	for _, c := range colNames {
		partColDef = append(partColDef, c+" "+colDef[c])
	}
	partPk := strings.Join(pk, ",")
	if partPk != "" {
		partPk = ", PRIMARY KEY (" + partPk + ")"
	}
	sql := fmt.Sprintf(template, partIfNotExists, tableName, strings.Join(partColDef, ","), partPk)
	_, err := sqlc.GetDB().Exec(sql)
	return err
}

// CreateIndex generates and executes "CREATE INDEX" statement
func CreateIndex(sqlc *prom.SqlConnect, tableName string, unique bool, cols []string) error {
	template := "CREATE INDEX idx_%s_%s on %s(%s)"
	templateUnique := "CREATE UNIQUE INDEX udx_%s_%s on %s(%s)"
	var sql string
	if unique {
		sql = fmt.Sprintf(templateUnique, tableName, strings.Join(cols, "_"), tableName, strings.Join(cols, ","))
	} else {
		sql = fmt.Sprintf(template, tableName, strings.Join(cols, "_"), tableName, strings.Join(cols, ","))
	}
	_, err := sqlc.GetDB().Exec(sql)
	return err
}
