package jobs

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err *error

const DNS = "root:admin@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"

type Jobs struct {
	gorm.Model
	jobid     string `json:"jobid"`
	jobstatus string `json:"jobstatus"`
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

	// Getting the job ID from the URL query parameters
	jobIDStr := c.Query("id")
	if jobIDStr == "" {
		c.Status(fiber.StatusBadRequest).Send("Missing job ID")
		return
	}

	// Converting the job ID to an integer
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		c.Status(fiber.StatusBadRequest).Send("Invalid job ID")
		return
	}

	// Look up the job by ID
	job, ok := jobs[jobID]
	if !ok {
		c.Status(fiber.StatusNotFound).Send("Job not found")
		return
	}

	// Return the status of the job
	if(job.Status=200)
	fmt.Println("Processed")
	else
	fmt.Println("Processing")

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
