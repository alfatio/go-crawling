package routes

import (
	"go-crawling/models"
	"time"

	"github.com/gin-gonic/gin"
)

// get kurs based on date on query
func GetKurs(c *gin.Context) {
	var args []interface{}

	var allKurs []models.Kurs

	symbol := c.Param("symbol")

	// check and assign startdate
	if d, err := time.Parse("2006-01-02", c.Query("startdate")); err == nil && c.Query("startdate") != "" {
		args = append(args, d.Format("2006-01-02"))
	} else if err != nil {
		c.JSON(400, gin.H{
			"message": "wrong date format",
		})

		c.Abort()
		return
	} else {
		args = append(args, "1111-01-01")
	}
	// check and assign enddate
	if d, err := time.Parse("2006-01-02", c.Query("enddate")); err == nil && c.Query("enddate") != "" {
		args = append(args, d.Format("2006-01-02"))
	} else if err != nil {
		c.JSON(400, gin.H{
			"message": "wrong date format",
		})

		c.Abort()
		return
	} else {
		args = append(args, "9999-12-31")
	}

	selectQuery := "select * from kurs where kurs_date >= ? and kurs_date <= ? "

	if symbol != "" {
		selectQuery += "and symbol = ?"
		args = append(args, symbol)
	}

	res, err := database.Query(selectQuery, args...)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var kurs models.Kurs

		err = res.Scan(
			&kurs.ID,
			&kurs.Symbol,
			&kurs.E_rate.Beli,
			&kurs.E_rate.Jual,
			&kurs.TT_counter.Beli,
			&kurs.TT_counter.Jual,
			&kurs.Bank_notes.Beli,
			&kurs.Bank_notes.Jual,
			&kurs.Date,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			panic(err)
		}

		allKurs = append(allKurs, kurs)
	}

	c.JSON(200, allKurs)

}

// insert kurs data from request body
func PostKurs(c *gin.Context) {
	var kursData models.Kurs

	err := c.BindJSON(&kursData)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "missing required field",
		})
		c.Abort()
		return
	}

	if len(kursData.Symbol) != 3 {
		c.JSON(400, gin.H{
			"message": "wrong symbol format",
		})
		c.Abort()
		return
	}

	if err := checkDate(kursData.Date); err != nil {
		c.JSON(400, gin.H{
			"message": "wrong date format",
		})

		c.Abort()
		return
	}

	row, err := database.Exec("insert ignore into kurs values (null, ?, ?, ?, ?, ?, ?, ?, ?)",
		kursData.Symbol,
		kursData.E_rate.Beli,
		kursData.E_rate.Jual,
		kursData.TT_counter.Beli,
		kursData.TT_counter.Jual,
		kursData.Bank_notes.Beli,
		kursData.Bank_notes.Jual,
		kursData.Date,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		panic(err)
	}

	if n, _ := row.RowsAffected(); n == 0 {
		c.JSON(200, gin.H{
			"message": "data already exist",
		})

		c.Abort()
		return
	}

	res := database.QueryRow("select last_insert_id()")

	if err := res.Scan(&kursData.ID); err == nil {
		c.JSON(200, kursData)
	} else {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		panic(err)
	}

}

// insert kurs data from request body
func PutKurs(c *gin.Context) {
	var kursData models.Kurs

	err := c.BindJSON(&kursData)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "missing required field",
		})
		c.Abort()
		return
	}

	if len(kursData.Symbol) != 3 {
		c.JSON(400, gin.H{
			"message": "wrong symbol format",
		})
		c.Abort()
		return
	}

	row, err := database.Exec(`
		update kurs set
			symbol = ?,
			e_rate_beli = ?, 
			e_rate_jual = ?,
			tt_counter_beli = ?,
			tt_counter_jual = ?,
			bank_notes_beli = ?,
			bank_notes_jual = ?,
			kurs_date = ?
		where symbol = ? and kurs_date = ?
	`,
		kursData.Symbol,
		kursData.E_rate.Beli,
		kursData.E_rate.Jual,
		kursData.TT_counter.Beli,
		kursData.TT_counter.Jual,
		kursData.Bank_notes.Beli,
		kursData.Bank_notes.Jual,
		kursData.Date,
		kursData.Symbol,
		kursData.Date,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		panic(err)
	}

	if n, _ := row.RowsAffected(); n == 0 {
		c.JSON(404, gin.H{
			"message": "data does not exist",
		})

		c.Abort()
		return
	}

	res := database.QueryRow("select * from kurs where symbol = ? and kurs_date = ?", kursData.Symbol, kursData.Date)

	var kurs models.Kurs

	err = res.Scan(
		&kurs.ID,
		&kurs.Symbol,
		&kurs.E_rate.Beli,
		&kurs.E_rate.Jual,
		&kurs.TT_counter.Beli,
		&kurs.TT_counter.Jual,
		&kurs.Bank_notes.Beli,
		&kurs.Bank_notes.Jual,
		&kurs.Date,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		panic(err)
	} else {
		c.JSON(200, kurs)
	}
}

// delete kurs based on date on params
func DeleteKurs(c *gin.Context) {
	date, err := time.Parse("2006-01-02", c.Param("date"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "wrong date format",
		})
	} else {
		_, err = database.Exec("delete from kurs where kurs_date = ?", date)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "ok",
			})
		}
	}
}

// change date string into time eg 2021-09-17 ->
func checkDate(d string) error {
	_, err := time.Parse("2006-01-02", d)
	if err != nil {
		return err
	} else {
		return nil
	}
}
