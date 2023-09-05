package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var settingsTableName = "settings"

func GetAppPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return exPath
}

func OpenDatabase() error {
	var err error

	db, err = sql.Open("sqlite3", GetAppPath()+"/sqlite-database.db")

	if err != nil {
		return err
	}

	return db.Ping()
}

func CreateTable() {
	runQuery(`DROP TABLE IF EXISTS ` + settingsTableName + `;`)
	runQuery(`CREATE TABLE ` + settingsTableName + ` (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "key" TEXT,
        "value" TEXT
      );`)
}

func IsSettingsTableCreated() bool {
	_, err := db.Query(`SELECT * FROM ` + settingsTableName)
	if err != nil {
		return false
	}
	return true
}

func FetchSettings() *Settings {
	rows, err := db.Query(`SELECT * FROM ` + settingsTableName)

	if err != nil {
		log.Println("Please init app first: ai init")
		log.Fatal(err)
	}

	settings := GetDefaultSettings()

	defer rows.Close()

	for rows.Next() {

		var id int
		var key string
		var value string

		err = rows.Scan(&id, &key, &value)

		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case "apiKey":
			decrypted, _ := DecryptMessage(GetEncryptionKey(), value)
			settings.ApiKey = decrypted
		case "apiUrl":
			settings.ApiURL = value
		case "model":
			settings.Model = value
		case "systemMessage":
			settings.SystemMessage = value
		case "maxTokens":
			settings.MaxTokens = StrToInt(value)
		case "temperature":
			settings.Temperature = StrToFloat(value)
		case "frequencyPenalty":
			settings.FrequencyPenalty = StrToFloat(value)
		case "presencePenalty":
			settings.PresencePenalty = StrToFloat(value)

		default:
			fmt.Printf(key, "not found")
		}

	}

	if len(settings.ApiKey) == 0 {
		InitOrDie()
	}

	return settings
}

func InsertKey(key string, value string) {
	insertNoteSQL := `INSERT INTO ` + settingsTableName + `(key, value) VALUES (?, ?)`
	statement, err := db.Prepare(insertNoteSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(key, value)
	if err != nil {
		log.Fatalln(err)
	}

}

func GetDefaultSettings() *Settings {

	return &Settings{
		ApiKey:           "",
		ApiURL:           "https://api.openai.com/v1/chat/completions",
		MaxTokens:        300,
		Temperature:      0.9,
		FrequencyPenalty: 0.9,
		PresencePenalty:  0.9,
		Model:            "gpt-4",
		SystemMessage:    "You are a senior software developer."}
}

func runQuery(query string) {
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}
