package databases

import (
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

func (sqlConn *SQLConnection) OpenSqlConnection() (*sql.DB, error) {
	defer sqlConn.mutex.Unlock()
	if sqlConn.db == nil {
		sqlConn.mutex.Lock()
		if sqlConn.db == nil {
			db, err := sql.Open("mysql", helpers.ReadEnv(helpers.MYSQLCONNECTIONSTRING))
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
	Begin()(*sql.Tx, error)
	Prepare(tx *sql.Tx, query string)(*sql.Stmt, error)
	Exec(stmt *sql.Stmt,args ...interface{})(sql.Result, error)
	LastInsertedID(result sql.Result)(int64, error)
	RowEffected(result sql.Result)(int64, error)
	Commit(tx *sql.Tx)error
	RollBack(tx *sql.Tx)error
}

type SQLClientImpl struct {
	sqlClient *sql.DB
}

func NewSqlClient(sql *sql.DB)SqlClient{
	return &SQLClientImpl{sqlClient:sql}
}

func (sql *SQLClientImpl) Begin()(*sql.Tx, error){
	return sql.sqlClient.Begin()
}

func (sql *SQLClientImpl) Prepare(tx *sql.Tx, query string)(*sql.Stmt, error){

	return tx.Prepare(query)
}

func (sql *SQLClientImpl) Exec(stmt *sql.Stmt,args ...interface{})(sql.Result, error){
	return stmt.Exec(args...)
}

func (sql *SQLClientImpl) LastInsertedID(result sql.Result)(int64, error){
	return result.LastInsertId()
}

func (sql *SQLClientImpl) RowEffected(result sql.Result)(int64, error){
	return result.RowsAffected()
}

func (sql *SQLClientImpl) Commit(tx *sql.Tx)error{
	return tx.Commit()
}

func (sql *SQLClientImpl) RollBack(tx *sql.Tx)error{
	return tx.Rollback()
}




