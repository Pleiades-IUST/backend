package drive

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Drive struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64
}

func (Drive) TableName() string {
	return "drive"
}

type Signal struct {
	ID              int64       `json:"id"`
	DriveID         int64       `json:"drive_id"`
	PlmnID          *string     `json:"plmn_id"`
	CellID          *string     `json:"cell_id"`
	Technology      *string     `json:"technology"`
	SignalStrength  *int32      `json:"signal_strength"`
	DownloadRate    *float64    `json:"download_rate"`
	UploadRate      *float64    `json:"upload_rate"`
	DnsLookupTime   *float64    `json:"dns_lookup_time"`
	Ping            *float64    `json:"ping"`
	SmsDeliveryTime *float64    `json:"sms_delivery_time"`
	RSRP            *int32      `json:"rsrp"`
	RSRQ            *int32      `json:"rsrq"`
	Longitude       *float64    `json:"longitude"`
	Latitude        *float64    `json:"latitude"`
	PCI             *string     `json:"pci"`
	TAC             *string     `json:"tac"`
	RecordTime      *CustomTime `json:"record_time"`
}

func (Signal) TableName() string {
	return "signal"
}

type CustomTime struct {
	time.Time
}

const customTimeLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes

	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(customTimeLayout, s, loc)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Value implements the driver.Valuer interface for database serialization
func (ct CustomTime) Value() (driver.Value, error) {
	if ct.IsZero() {
		return nil, nil
	}
	return ct.Time.Format(customTimeLayout), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case []byte:
		return ct.parseString(string(v))
	case string:
		return ct.parseString(v)
	default:
		return fmt.Errorf("cannot scan type %T into CustomTime", value)
	}
}

func (ct *CustomTime) parseString(s string) error {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation(customTimeLayout, s, loc)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}
