package asakamiplugins

import (
	"database/sql"
	"fmt"
	"log"
)

type CURRENCY struct {
	Topic   string
	Tooltip string
	itemid  int
	value   float64
}

type economystruct struct {
	pluginversion  string
	pluginhelp     string
	GetCurrency    func(itemid int) CURRENCY
	InsertCurrency func(currency CURRENCY)
	UpdateCurrency func(currency CURRENCY)
}

var economyv = economystruct{
	pluginversion:  "0.1.0",
	pluginhelp:     "- 浅上经济插件",
	GetCurrency:    GetCurrency,
	InsertCurrency: InsertCurrency,
	UpdateCurrency: UpdateCurrency,
}

// Geteconomy 导出变量
func Geteconomy() economystruct {
	return economyv
}

func createCurrencyTable(db *sql.DB) {
	// Create the currency table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS currency (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			topic TEXT,
			tooltip TEXT,
			itemid INTEGER,
			value REAL
		);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Print(err)
	}
}

func GetCurrency(itemid int) CURRENCY {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/currency.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Get currency values
	getCurrencySQL := `
		SELECT
			topic,
			tooltip,
			itemid,
			value
		FROM currency
		WHERE itemid = ?
	`

	row := db.QueryRow(getCurrencySQL, itemid)

	var currency CURRENCY
	err = row.Scan(&currency.Topic, &currency.Tooltip, &currency.itemid, &currency.value)
	if err != nil {
		log.Print(err)
	}

	return currency
}

func InsertCurrency(currency CURRENCY) {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/currency.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Insert currency values
	insertCurrencySQL := `
		INSERT INTO currency(
			topic,
			tooltip,
			itemid,
			value
		) VALUES (?, ?, ?, ?)
	`

	_, err = db.Exec(insertCurrencySQL, currency.Topic, currency.Tooltip, currency.itemid, currency.value)
	if err != nil {
		log.Print(err)
	}
}

func UpdateCurrency(currency CURRENCY) {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/currency.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Update currency values
	updateCurrencySQL := `
		UPDATE currency SET
			topic = ?,
			tooltip = ?,
			itemid = ?,
			value = ?
		WHERE itemid = ?
	`

	_, err = db.Exec(updateCurrencySQL, currency.Topic, currency.Tooltip, currency.itemid, currency.value, currency.itemid)
	if err != nil {
		log.Print(err)
	}
}

func init() {
	// Open the database connection
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/currency.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Create the currency table
	createCurrencyTable(db)

	//如果货币表为空则插入货币
	c := GetCurrency(1001)
	if c.Topic == "" {
		for _, currency := range currencylist {
			InsertCurrency(currency)
		}
		fmt.Println("Currency table created successfully")
	}

}

var currencylist = []CURRENCY{
	{
		Topic:   "稻荷币",
		Tooltip: "通用货币",
		itemid:  1001,
		value:   1,
	}, {
		Topic:   "钓鱼币",
		Tooltip: "mcfish插件使用货币",
		itemid:  1002,
		value:   0.01,
	}, {
		Topic:   "円",
		Tooltip: "日元",
		itemid:  1003,
		value:   5,
	},
}
