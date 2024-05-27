package asakamiplugins

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/FloatTech/floatbox/file"
)

var itemlist = []BackpackItem{
	{0, "EXAMPLE", "ABC", 1, "名称", "提示", "效果", 0, "标签"},
	{0, "CURRENCY", "DHC", 1001, "稻荷币", "通用货币", "", 0, ""},
	{0, "CURRENCY", "DYC", 1002, "钓鱼币", "mcfish插件使用货币", "", 0, ""},
	{0, "CURRENCY", "JPY", 1003, "円", "日元", "", 0, ""},
	{0, "CURRENCY", "CNY", 1004, "元", "人民币", "", 0, ""},
	{0, "CURRENCY", "USD", 1005, "美元", "美元", "", 0, ""},
}

type BackpackItem struct {
	ID       int
	ItemType string
	SubType  string
	itemid   int
	Topic    string
	Tooltip  string
	Buff     string
	Quantity float64
	Tag      string
}

type backpackstruct struct {
	pluginversion string
	pluginhelp    string
	InsertItem    func(item BackpackItem, uid int64)
	GetItem       func(itemid int, uid int64) (BackpackItem, error)
	GetItems      func(uid int64) []BackpackItem
	UpdateItem    func(item BackpackItem, uid int64)
	DeleteItem    func(id int, uid int64)
	GetItemlist   func() []BackpackItem
	init          func(uid int64)
}

var backpackv = backpackstruct{
	pluginversion: "0.1.0",
	pluginhelp:    "- 浅上背包插件",
	InsertItem:    insertItem,
	GetItem:       getItem,
	GetItems:      getItems,
	UpdateItem:    updateItem,
	DeleteItem:    deleteItem,
	GetItemlist: func() []BackpackItem {
		return itemlist
	},
	init: func(uid int64) {
		createpersonbackpack(uid)
	},
}

// Getbackpack 导出变量
func Getbackpack() backpackstruct {
	return backpackv
}

// InsertItem 插入物品
func insertItem(item BackpackItem, uid int64) {
	createpersonbackpack(uid)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	insertItemSQL := `INSERT INTO backpack` + fmt.Sprintf("%d", uid) + ` (itemType, subType,itemid, topic, tooltip, buff, quantity, tag) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = db.Exec(insertItemSQL, item.ItemType, item.SubType, item.itemid, item.Topic, item.Tooltip, item.Buff, item.Quantity, item.Tag)
	if err != nil {
		log.Print(err)
	}
}

// GetItem 获取物品
func getItem(itemid int, uid int64) (BackpackItem, error) {
	createpersonbackpack(uid)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	getItemSQL := `SELECT * FROM backpack` + fmt.Sprintf("%d", uid) + ` WHERE itemid = ?;`

	var item BackpackItem

	err = db.QueryRow(getItemSQL, itemid).Scan(&item.ID, &item.ItemType, &item.SubType, &item.itemid, &item.Topic, &item.Tooltip, &item.Buff, &item.Quantity, &item.Tag)
	if err != nil {
		log.Print(err)
		item, err = getitemfromlist(itemid)
		if err == nil {
			insertItem(item, uid)
			log.Print("created item")
		}
	}

	return item, err
}

// 从itemlist中由itemid获取item
func getitemfromlist(itemid int) (BackpackItem, error) {
	for _, item := range itemlist {
		if item.itemid == itemid {
			return item, nil
		}
	}
	return itemlist[0], errors.New("item not found")
}

// GetItems 获取所有物品
func getItems(uid int64) []BackpackItem {
	createpersonbackpack(uid)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	getItemsSQL := `
		SELECT * FROM backpack` + fmt.Sprintf("%d", uid) + `;
	`

	rows, err := db.Query(getItemsSQL)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	var items []BackpackItem

	for rows.Next() {
		var item BackpackItem
		err = rows.Scan(&item.ID, &item.ItemType, &item.SubType, &item.itemid, &item.Topic, &item.Tooltip, &item.Buff, &item.Quantity, &item.Tag)
		if err != nil {
			log.Print(err)
		}
		items = append(items, item)
	}

	return items
}

// UpdateItem 更新物品
func updateItem(item BackpackItem, uid int64) {
	createpersonbackpack(uid)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	updateItemSQL := `
		UPDATE backpack` + fmt.Sprintf("%d", uid) + ` SET itemType = ?, subType = ?, itemid = ?, topic = ?, tooltip = ?, buff = ?, quantity = ?, tag = ? WHERE id = ?;
	`

	_, err = db.Exec(updateItemSQL, item.ItemType, item.SubType, item.itemid, item.Topic, item.Tooltip, item.Buff, item.Quantity, item.Tag, item.ID)
	if err != nil {
		log.Print(err)
	}
}

// DeleteItem 删除物品
func deleteItem(id int, uid int64) {
	createpersonbackpack(uid)
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	deleteItemSQL := `
		DELETE FROM backpack` + fmt.Sprintf("%d", uid) + ` WHERE id = ?;
	`

	_, err = db.Exec(deleteItemSQL, id)
	if err != nil {
		log.Print(err)
	}
}

func init() {
	if file.IsNotExist("data/asakamiplugins/db") {
		err := os.MkdirAll("data/asakamiplugins/db", 0755)
		if err != nil {
			panic(err)
		}
	}
	/*
		//获取所有表
		db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
		if err != nil {
			log.Print(err)
		}
		defer db.Close()
		//获取所有表
		getTablesSQL := `
			SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;
		`

		rows, err := db.Query(getTablesSQL)
		if err != nil {
			log.Print(err)
		}

		var tables []string

		for rows.Next() {
			var table string
			err = rows.Scan(&table)
			if err != nil {
				log.Print(err)
			}
			tables = append(tables, table)
		}

		//获取每个表itemid=1003的物品
		for _, table := range tables {

			getItemSQL := `SELECT * FROM ` + table + ` WHERE itemid = ?;`

			var item BackpackItem

			err = db.QueryRow(getItemSQL, 1003).Scan(&item.ID, &item.ItemType, &item.SubType, &item.itemid, &item.Topic, &item.Tooltip, &item.Buff, &item.Quantity, &item.Tag)
			if err != nil {
				log.Print(err)
				item = getitemfromlist(1003)
				insertItem(item, 0)
				log.Print("created item")
			}
		}

		//将quantity除以25
		for _, table := range tables {

			getItemSQL := `SELECT * FROM ` + table + ` WHERE itemid = ?;`

			var item BackpackItem

			err = db.QueryRow(getItemSQL, 1003).Scan(&item.ID, &item.ItemType, &item.SubType, &item.itemid, &item.Topic, &item.Tooltip, &item.Buff, &item.Quantity, &item.Tag)
			if err != nil {
				log.Print(err)
				item = getitemfromlist(1003)
				insertItem(item, 0)
				log.Print("created item")
			}
			item.Quantity = item.Quantity / 25
			//更新物品
			updateItemSQL := `
			UPDATE ` + table + ` SET itemType = ?, subType = ?, itemid = ?, topic = ?, tooltip = ?, buff = ?, quantity = ?, tag = ? WHERE id = ?;
		`

			_, err = db.Exec(updateItemSQL, item.ItemType, item.SubType, item.itemid, item.Topic, item.Tooltip, item.Buff, item.Quantity, item.Tag, item.ID)
			if err != nil {
				log.Print(err)
			}
		}
	*/
}

func createpersonbackpack(uid int64) {
	db, err := sql.Open("sqlite3", "data/asakamiplugins/db/backpack.db")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS backpack` + fmt.Sprintf("%d", uid) + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemType TEXT,
		subType TEXT,
		itemid INTEGER,
		topic TEXT,
		tooltip TEXT,
		buff TEXT,
		quantity REAL,
		tag TEXT
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Print(err)
	}
}
