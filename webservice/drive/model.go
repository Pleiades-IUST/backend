package drive

import "time"

type Drive struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (Drive) TableName() string {
	return "drive"
}

type Signal struct {
	ID         int64  `json:"id"`
	DriveID    int64  `json:"drive_id"`
	Technology string `json:"technology"`
	Strength   int32  `json:"strength"`
	RSRP       int32  `json:"rsrp"`
	RSRQ       int32  `json:"rsrq"`
}

func (Signal) TableName() string {
	return "signal"
}
