package jobs

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err *error

const DNS = "root:admin@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"

type ConnectionDetails struct {
	user     string
	password string
	host     string
	port     string
	database string
}

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to Database")
	}
	DB.AutoMigrate(&Jobs{})
}

func Get(c *fiber.Ctx) error {

	// Get the id parameter from the request
	id := c.Params("id")

	// Querying the database to see if the id matches the lastid of any row
	var lastid string
	err := db.QueryRow("SELECT lastid FROM table WHERE lastid = ?", id).Scan(&lastid)
	if err != nil {
		// If the id doesn't match the lastid of any row, return "Processing"
		if err == sql.ErrNoRows {
			c.Send("Processing")
			return
		}
		// If there was an error executing the query, log the error and return an error message
		log.Println(err)
		c.Send("Error executing query")
		return
	}

	// If the id matches the lastid of a row, return "Processed"
	c.Send("Processed")
}

func Post(c *fiber.Ctx) error {
	csvFile, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Read the data from the CSV file
	reader := csv.NewReader(csvFile)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Insert the data into the MySQL database
		// Insert a new row into the table
		res, err := db.Exec("INSERT INTO  (Name, Age, Roll No., Department) VALUES (?, ?, ?, ?)", row[0], row[1], row[2], row[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get the ID of the last inserted row
		//This is mainly an SQL function named LastInsertId which is kind of an address for the row inserted, and is unique
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Last inserted ID:", lastInsertId)
	}

	c.Send("Data imported successfully")
}
