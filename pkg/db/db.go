package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	err    error
	logger = log.New(os.Stdout, "[database]", log.Lshortfile)
)

func ConnectDb() error {
	db, err = sql.Open("mysql", "root:@/github-network")
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	logger.Println("connected!")
	return nil
}

func ExecuteQuery(query string, args ...interface{}) error {
	if _, err := db.Exec(query, args...); err != nil {
		logger.Println("failed in executing query!")
		logger.Println(err)
		return err
	}
	return nil
}

func QueryUserLogins() ([]string, error) {
	query := `
		SELECT user_login FROM pull_requests
		UNION
		SELECt user_login FROM issues;
	`
	rows, err := db.Query(query)
	if err != nil {
		logger.Println("failed in executing query!")
		logger.Println(err)
		return nil, err
	}

	defer rows.Close()
	var user_logins []string
	for rows.Next() {
		var user_login string
		if err := rows.Scan(&user_login); err != nil {
			logger.Println("failed in scanning row!")
			return nil, err
		}
		user_logins = append(user_logins, user_login)
	}
	if err := rows.Err(); err != nil {
		logger.Println("sql error!")
		return nil, err
	}

	return user_logins, nil
}
