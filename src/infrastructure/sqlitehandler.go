package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/anhbkpro/go-cleanarchitecture/src/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	Conn *sql.DB
}

// DbHandler implementation
func (handler *SqliteHandler) Execute(statement string) {
	handler.Conn.Exec(statement)
}

func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}

	row := new(SqliteRow)
	row.Rows = rows

	return row
}

// interfaces.Row implementation
type SqliteRow struct {
	Rows *sql.Rows
}

func (r SqliteRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

func NewSqliteHandler(dbFilename string) *SqliteHandler {
	conn, _ := sql.Open("sqlite3", dbFilename)
	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = conn
	return sqliteHandler
}
