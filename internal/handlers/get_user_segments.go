package handlers

import (
	"fmt"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetUserSegments(c *gin.Context) {
	timeToDb := time.Now()
	var requestData struct {
		UserID int `json:"user_id"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	db, errDb := databases.CreateDatabaseConnect()
	if errDb != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "couldn't open database"})
		return
	}
	defer db.Close()

	query := `
		SELECT ids.segment_name
		FROM users_segments us
		INNER JOIN id_name_segments ids ON us.segment_id = ids.segment_id
		WHERE us.user_id = $1 AND us.endtime > $2
	`
	fmt.Println(db.Stats(), "OK")
	data, errGet := db.Query(query, requestData.UserID, timeToDb)
	if errGet != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errGet})
		return
	}
	defer data.Close()

	var segments []string
	for data.Next() {
		var segmentName string
		if err := data.Scan(&segmentName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		segments = append(segments, segmentName)
	}

	if err := data.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		return
	}
	c.JSON(http.StatusOK, segments)
}
