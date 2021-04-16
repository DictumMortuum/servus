package weight

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type WeightRow struct {
	Id               int64     `db:"id"`
	DateRaw          string    `form:"date"`
	Date             time.Time `db:"date"`
	Weight           float64   `db:"weight" form:"weight"`
	Bmi              float64   `db:"bmi" form:"bmi"`
	Bmr              int       `db:"bmr" form:"bmr"`
	Impedance        int       `db:"impedance" form:"impedance"`
	Fat              float64   `db:"fat" form:"fat"`
	FatMass          float64   `db:"fat_mass" form:"fat_mass"`
	PredictedWeight  float64   `db:"predicted_weight" form:"predicted_weight"`
	PredictedFatMass float64   `db:"predicted_fat_mass" form:"predicted_fat_mass"`
	FatToLose        float64   `db:"fat_to_lose" form:"fat_to_lose"`
}

func CreateWeight(db *sqlx.DB, data WeightRow) error {
	sql := `
	insert into tweight (
		date,
		weight,
		bmi,
		bmr,
		impedance,
		fat,
		fat_mass,
		predicted_weight,
		predicted_fat_mass,
		fat_to_lose
	) values (
		:date,
		:weight,
		:bmi,
		:bmr,
		:impedance,
		:fat,
		:fat_mass,
		:predicted_weight,
		:predicted_fat_mass,
		:fat_to_lose
	)`

	_, err := db.NamedExec(sql, &data)
	if err != nil {
		return err
	}

	return nil
}
