package helpers

import (
	"bufio"
	"encoding/json"
	"mime/multipart"
	"os"
)

const (
	MYSQLCONNECTIONSTRING = "MysqlConnectionStr"
	DATABASENAME          = "inventory"
	HELPER_VERSION_V1     = "v1"
	MigrationFileLocation = "MIGRATIONFILE"
)

type HelperBase interface {
	SetEnv()
	ReadEnv(envType string) string
	StreamFile(file multipart.File) (*json.Decoder, error)
}

type Helper struct {
}

func NewHelper(version string) HelperBase {
	switch version {
	case HELPER_VERSION_V1:
		return &Helper{}
	default:
		return &Helper{}
	}
}
func (helper *Helper) SetEnv() {
	// Containers env set
	//os.Setenv("MYSQL_CONN", "root:root@tcp(localhost:3308)/inventory?")
	//os.Setenv("MIGRATIONFILE", "file://")

	os.Setenv("MYSQL_CONN", "root:Siar@123@tcp(localhost:3306)/inventory?")
	os.Setenv("MIGRATIONFILE", "file://pkg/migration_scripts")
}

func (helper *Helper) ReadEnv(envType string) string {
	switch envType {
	case MYSQLCONNECTIONSTRING:
		return os.Getenv("MYSQL_CONN")
	case MigrationFileLocation:
		return os.Getenv("MIGRATIONFILE")
	default:
		return ""
	}
}

func (helper *Helper) StreamFile(file multipart.File) (*json.Decoder, error) {
	r := bufio.NewReader(file)
	d := json.NewDecoder(r)
	return d, nil
}
