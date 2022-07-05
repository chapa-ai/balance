package db

import (
	"balance/pkg/models"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "123"
//	dbname   = "postgres"
//)

func InitDB() (*sql.DB, error) {
	if db == nil {
		conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
		d, err := sql.Open("postgres", conn)
		if err != nil {
			return nil, err
		}
		db = d
	}
	return db, nil
}

func MigrateDb() error {
	db, err := InitDB()
	if err != nil {
		return err
	}

	migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "migrations")
	if err != nil {
		return err
	}
	return migrator.Migrate()
}

func SelectBalance(userId string) (float64, error) {
	db, err := InitDB()
	if err != nil {
		return 0, err
	}
	stmt := fmt.Sprintf(`SELECT balance  FROM "balance" WHERE "userid"=$1`)

	var balance float64

	_ = db.QueryRow(stmt, userId).Scan(&balance)

	return balance, nil
}

func ChangeTable(data *models.Balance, total float64) (float64, error) {
	conn, err := InitDB()
	if err != nil {
		return 0, err
	}

	var balance float64

	err = conn.QueryRow(`INSERT INTO "balance"("userid", "balance") VALUES($1, $2) ON CONFLICT (userid) 
    DO UPDATE SET balance = $3 RETURNING "balance"`, data.UserId, total, total).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func UpdateTables(data *models.Balance) ([]*models.Balance, error) {

	conn, err := InitDB()
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	list := []*models.Balance{}

	balanceUser1, err := SelectBalance(data.UserId)
	if err != nil {
		return nil, err
	}

	totalUser1 := balanceUser1 - data.Sum

	row, err := conn.Query(`INSERT INTO "balance"("userid", "balance") VALUES($1, $2) ON CONFLICT (userid)
	DO UPDATE SET balance = $2 RETURNING "userid", "balance"`, data.UserId, totalUser1)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		obj := models.Balance{}
		if err := row.Scan(&obj.UserId, &obj.Balance); err != nil {
			return nil, err
		}
		list = append(list, &obj)
	}

	balanceUser2, err := SelectBalance(data.UserId2)
	if err != nil {
		return nil, err
	}

	totalUser2 := balanceUser2 + data.Sum

	row2, err := conn.Query(`INSERT INTO "balance"("userid", "balance") VALUES($1, $2) ON CONFLICT (userid) 
    DO UPDATE SET balance = $2 RETURNING "userid", "balance"`, data.UserId2, totalUser2)
	if err != nil {
		return nil, err
	}

	for row2.Next() {
		obj2 := models.Balance{}
		if err := row2.Scan(&obj2.UserId, &obj2.Balance); err != nil {
			return nil, err
		}
		list = append(list, &obj2)
	}

	return list, nil
}

func SendMoney(data *models.Balance) (float64, error) {
	var balance float64

	balance, err := SelectBalance(data.UserId)
	if err != nil {
		return 0, err
	}
	total := balance + data.Sum
	if total < 0 {
		return 0, fmt.Errorf("You haven't got enough money for making this operation", err)
	}

	balance, err = ChangeTable(data, total)
	if err != nil {
		return 0, err
	}

	return balance, nil

}
