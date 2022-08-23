package routes

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"

	"go-crawling/db"
	"go-crawling/models"
)

var database = db.DB()

func Indexing(c *gin.Context) {
	var allKurs []models.Kurs

	co := colly.NewCollector()
	var dateKurs string

	co.OnHTML("tr[code]", func(h *colly.HTMLElement) {
		code := h.Attr("code")

		var kurs models.Kurs

		// asign symbol
		kurs.Symbol = code

		// assign e_rate jual & beli
		e_rate_beli := h.ChildText("p[rate-type=ERate-buy]")

		kurs.E_rate.Beli = changeDecimal(e_rate_beli)

		e_rate_jual := h.ChildText("p[rate-type=ERate-sell]")

		kurs.E_rate.Jual = changeDecimal(e_rate_jual)

		// assign tt_counter jual & beli

		tt_counter_beli := h.ChildText("p[rate-type=TT-buy]")

		kurs.TT_counter.Beli = changeDecimal(tt_counter_beli)

		tt_counter_jual := h.ChildText("p[rate-type=TT-sell]")

		kurs.TT_counter.Jual = changeDecimal(tt_counter_jual)

		// assign bank_notes jual & beli

		bn_beli := h.ChildText("p[rate-type=BN-buy]")

		kurs.Bank_notes.Beli = changeDecimal(bn_beli)

		bn_jual := h.ChildText("p[rate-type=BN-sell]")

		kurs.Bank_notes.Jual = changeDecimal(bn_jual)

		allKurs = append(allKurs, kurs)
	})

	co.OnHTML(".desc-ref-kurs.refresh-date", func(h *colly.HTMLElement) {
		d := h.Text
		d = strings.Join(strings.Split(d, " ")[0:3], " ")

		dateTaken, err := time.Parse("02 January 2006", d)

		if err != nil {
			panic(err)
		}

		result := dateTaken.Format("2006-01-02")

		dateKurs = result
	})

	//start scrapping
	co.Visit("https://www.bca.co.id/id/informasi/kurs")

	// check if kurs with date already exist on db
	countData := 0

	err := database.QueryRow(`select count(*) from kurs where kurs_date = ? limit 1`, dateKurs).Scan(&countData)

	if err != nil {
		panic(err)
	}

	if countData != 0 {
		c.JSON(200, gin.H{
			"message": "today kurs already indexed",
		})

	} else {

		// insert all kurs data
		insertQuery := `insert into kurs values %s`
		var valuesStrings []string
		var valuesArgs []interface{}

		for _, el := range allKurs {
			valuesStrings = append(valuesStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")

			valuesArgs = append(valuesArgs, sql.NullInt16{})
			valuesArgs = append(valuesArgs, el.Symbol)
			valuesArgs = append(valuesArgs, el.E_rate.Beli)
			valuesArgs = append(valuesArgs, el.E_rate.Jual)
			valuesArgs = append(valuesArgs, el.TT_counter.Beli)
			valuesArgs = append(valuesArgs, el.TT_counter.Jual)
			valuesArgs = append(valuesArgs, el.Bank_notes.Beli)
			valuesArgs = append(valuesArgs, el.Bank_notes.Jual)
			valuesArgs = append(valuesArgs, dateKurs)
		}

		insertQuery = fmt.Sprintf(insertQuery, strings.Join(valuesStrings, ","))

		_, err := database.Exec(insertQuery, valuesArgs...)

		if err != nil {
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			panic(err)
		} else {
			c.JSON(200, gin.H{
				"message": "indexing finish",
			})
		}

	}

}

// change string currency into float64 eg: 14.501,00 -> 14501.00
func changeDecimal(n string) float64 {
	str := strings.Replace(n, ".", "", -1)
	str = strings.Replace(str, ",", ".", -1)

	flt, _ := strconv.ParseFloat(str, 64)
	return flt
}
