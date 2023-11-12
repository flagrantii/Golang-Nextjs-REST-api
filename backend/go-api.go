package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
type Attractions struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Detail     string `db:"detail" json:"detail"`
	Coverimage string `db:"coverimage" json:"coverimage"`
}

func getAttractions(c *gin.Context) {
	var attractions []Attractions
	rows, err := db.Query("select id, name, detail, coverimage from attractions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var a Attractions
		err := rows.Scan(&a.Id, &a.Name, &a.Detail, &a.Coverimage)
		if err != nil {
			log.Fatal(err)
		}
		attractions = append(attractions, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, attractions)
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/go_restapi")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/attractions", func(c *gin.Context) {
		getAttractions(c)
	})
	router.GET("/attractions/:id", func(c *gin.Context) {
		var attraction Attractions
		id := c.Param("id")
		row := db.QueryRow("select id, name, detail, coverimage from attractions where id = ?", id)
		err := row.Scan(&attraction.Id, &attraction.Name, &attraction.Detail, &attraction.Coverimage)
		if err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Attraction not found"})
				return
			}
			log.Fatal(err)
		}
		c.IndentedJSON(http.StatusOK, attraction)
	})

	router.POST("/attractions", func(c *gin.Context) {
		var newAttraction Attractions
		if err := c.BindJSON(&newAttraction); err != nil {
			return
		}
		stmt, err := db.Prepare("insert into attractions (name, detail, coverimage) values(?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(newAttraction.Name, newAttraction.Detail, newAttraction.Coverimage)
		if err != nil {
			log.Fatal(err)
		}
		c.IndentedJSON(http.StatusCreated, newAttraction)
	})

	router.PUT("/attractions/:id", func(c *gin.Context) {
		var newAttraction Attractions
		if err := c.BindJSON(&newAttraction); err != nil {
			return
		}
		id := c.Param("id")
		stmt, err := db.Prepare("update attractions set name=?, detail=?, coverimage=? where id=?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(newAttraction.Name, newAttraction.Detail, newAttraction.Coverimage, id)
		if err != nil {
			log.Fatal(err)
		}
		c.IndentedJSON(http.StatusOK, newAttraction)
	})

	router.DELETE("/attractions/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("delete from attractions where id=?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(id)
		if err != nil {
			log.Fatal(err)
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Attraction deleted"})
	})

	router.Run("localhost:8080")

}
