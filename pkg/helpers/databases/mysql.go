package databases

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"sync"
)

var connection *SQLConnection

type SQLConnection struct {
	db    *sql.DB
	mutex sync.Mutex
}

func NewSqlConnection() SQLConnection {
	return SQLConnection{}
}

func (sqlConn *SQLConnection) OpenSqlConnection(helper helpers.HelperBase) (*sql.DB, error) {
	defer sqlConn.mutex.Unlock()
	if sqlConn.db == nil {
		sqlConn.mutex.Lock()
		if sqlConn.db == nil {
			db, err := sql.Open("mysql", helper.ReadEnv(helpers.MYSQLCONNECTIONSTRING))
			if err != nil {
				return nil, err
			}
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(10)
			db.SetConnMaxLifetime(60)
			sqlConn.db = db
			connection = sqlConn
			return sqlConn.db, nil
		}
	}
	return sqlConn.db, nil
}

func GetSqlConnection() *sql.DB {
	return connection.db
}

type SqlClient interface {
	Begin() (*sql.Tx, error)
	BeginWithContext(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error)
	Prepare(tx *sql.Tx, query string) (*sql.Stmt, error)
	PrepareWithContext(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error)
	Exec(stmt *sql.Stmt, args ...interface{}) (sql.Result, error)
	ExecWithContext(ctx context.Context, stmt *sql.Stmt, args ...interface{}) (sql.Result, error)
	LastInsertedID(result sql.Result) (int64, error)
	RowEffected(result sql.Result) (int64, error)
	Commit(tx *sql.Tx) error
	RollBack(tx *sql.Tx) error
	GetIsolationLevel(isolationLevel int) sql.IsolationLevel
	BuildOptions(readOlny bool, isolationLevel sql.IsolationLevel) sql.TxOptions
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryWithContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowWithContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type SQLClientImpl struct {
	sqlClient *sql.DB
}

func NewSqlClient(sql *sql.DB) SqlClient {
	return &SQLClientImpl{sqlClient: sql}
}

func (sql *SQLClientImpl) Begin() (*sql.Tx, error) {
	return sql.sqlClient.Begin()
}

func (sql *SQLClientImpl) BeginWithContext(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	return sql.sqlClient.BeginTx(ctx, options)
}

func (sql *SQLClientImpl) Prepare(tx *sql.Tx, query string) (*sql.Stmt, error) {
	return tx.Prepare(query)
}

func (sql *SQLClientImpl) PrepareWithContext(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error) {
	return tx.PrepareContext(ctx, query)
}

func (sql *SQLClientImpl) Exec(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	return stmt.Exec(args...)
}

func (sql *SQLClientImpl) ExecWithContext(ctx context.Context, stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	return stmt.ExecContext(ctx, args...)
}

func (sql *SQLClientImpl) LastInsertedID(result sql.Result) (int64, error) {
	return result.LastInsertId()
}

func (sql *SQLClientImpl) RowEffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

func (sql *SQLClientImpl) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (sql *SQLClientImpl) RollBack(tx *sql.Tx) error {
	return tx.Rollback()
}

func (sq *SQLClientImpl) BuildOptions(readOlny bool, isolationLevel sql.IsolationLevel) sql.TxOptions {
	return sql.TxOptions{ReadOnly: readOlny, Isolation: isolationLevel}

}

func (sql *SQLClientImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := sql.sqlClient.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (sql *SQLClientImpl) QueryWithContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := sql.sqlClient.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (sql *SQLClientImpl) QueryRow(query string, args ...interface{}) *sql.Row {
	return sql.sqlClient.QueryRow(query, args...)
}

func (sql *SQLClientImpl) QueryRowWithContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return sql.sqlClient.QueryRowContext(ctx, query, args...)
}

func (sq *SQLClientImpl) GetIsolationLevel(isolationLevel int) sql.IsolationLevel {
	switch isolationLevel {
	case 1:
		return sql.LevelReadCommitted
	case 2:
		return sql.LevelReadUncommitted
	case 3:
		return sql.LevelWriteCommitted
	case 4:
		return sql.LevelRepeatableRead
	case 5:
		return sql.LevelSnapshot
	case 6:
		return sql.LevelSerializable
	case 7:
		return sql.LevelLinearizable
	default:
		return sql.LevelDefault
	}
}
