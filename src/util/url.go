package util

import "fmt"

func FormatUrl(addr string, port int) (url string) {
	url = fmt.Sprintf("%s:%d", addr, port)
	return
}

func FormatDatabaseDSN(address string, port int, username, password, database string) (dsn string) {
	dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, address, port, database)
	return dsn
}
func GetDefaultDSN() (dsn string) {
	config := GetConfig().Database
	dsn = FormatDatabaseDSN(config.Address, config.Port, config.Username, config.Password, config.DataBase)
	return
}
