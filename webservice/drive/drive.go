package drive

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Pleiades-IUST/backend/utils/dbutil"
	"github.com/Pleiades-IUST/backend/utils/ginutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateDrive(ctx *gin.Context) {
	testDriveData := struct {
		Drive   *Drive
		Signals []*Signal
	}{}

	userID := ginutil.GetUserID(ctx)

	err := ctx.ShouldBindJSON(&testDriveData)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if testDriveData.Drive == nil {
		testDriveData.Drive = &Drive{}
	}

	if len(testDriveData.Drive.Name) == 0 {
		testDriveData.Drive.Name = GenerateRandomString(10)
	}

	testDriveData.Drive.UserID = userID

	tx := dbutil.GormDB(ctx.Request.Context())

	err = tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&(testDriveData.Drive)).Error
		if err != nil {
			return err
		}

		for _, s := range testDriveData.Signals {
			s.DriveID = testDriveData.Drive.ID
		}

		err = tx.Create(testDriveData.Signals).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func FetchAllDrives(ctx *gin.Context) {
	tx := dbutil.GormDB(ctx)

	userID := ginutil.GetUserID(ctx)

	drives := []*Drive{}

	err := tx.Table("drive").Where("user_id = ?", userID).Scan(&drives).Error
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, drives)
}

func FetchSignals(ctx *gin.Context) {
	tx := dbutil.GormDB(ctx)

	userID := ginutil.GetUserID(ctx)

	request := struct {
		DriveID int64 `json:"drive_id"`
	}{}

	err := ctx.ShouldBindBodyWithJSON(&request)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	driveID := request.DriveID

	signals := []Signal{}

	err = tx.Table("signal AS s").
		Joins("JOIN drive AS d ON s.drive_id = d.id").
		Where("s.drive_id = ? AND d.user_id = ?", driveID, userID).
		Select("s.*").
		Scan(&signals).Error
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	result := struct {
		Signals []Signal `json:"signals"`
		Csv     string   `json:"csv"`
	}{
		Signals: signals,
		Csv:     fmt.Sprintf("http://localhost:8080/drive/csv?drive_id=%d", driveID),
	}

	ctx.JSON(http.StatusOK, result)
}

func GetCSV(ctx *gin.Context) {
	ctx.Header("Content-Disposition", "attachment; filename=data.csv")
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Transfer-Encoding", "binary")

	tx := dbutil.GormDB(ctx)

	driveIDstr := ctx.Query("drive_id")
	if driveIDstr == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	driveID, err := strconv.ParseInt(driveIDstr, 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	signals := []*Signal{}

	err = tx.Table("signal AS s").
		Joins("JOIN drive AS d ON s.drive_id = d.id").
		Where("d.id = ?", driveID).
		Select("s.*").
		Scan(&signals).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Create a CSV writer writing directly to the response
	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// Your CSV data
	data := convertSignalsToCSV(signals)

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			ctx.String(http.StatusInternalServerError, "Error writing CSV: %v", err)
			return
		}
	}
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomString)
}

func convertSignalsToCSV(signals []*Signal) [][]string {
	result := [][]string{}

	headers := []string{
		"plmn_id",
		"cell_id",
		"technology",
		"signal_strength",
		"dowload_rate",
		"upload_rate",
		"dns_lookup_time",
		"ping",
		"sms_delivery_time",
		"rsrp",
		"rsrq",
		"longitude",
		"latitude",
		"pci",
		"tac",
		"record_time",
		"lac",
		"rac",
		"frequency_band",
		"arfcn",
		"frequency",
		"rscp",
		"ecn0",
		"rxlev",
	}

	result = append(result, headers)

	for _, signal := range signals {
		newRow := []string{}

		newRow = append(newRow, nullStringToCSV(signal.PlmnID))

		newRow = append(newRow, nullStringToCSV(signal.CellID))

		newRow = append(newRow, nullStringToCSV(signal.Technology))

		newRow = append(newRow, nullInt32ToCsv(signal.SignalStrength))

		newRow = append(newRow, nullFloat64ToCSV(signal.DownloadRate, 3))

		newRow = append(newRow, nullFloat64ToCSV(signal.UploadRate, 3))

		newRow = append(newRow, nullFloat64ToCSV(signal.DnsLookupTime, 3))

		newRow = append(newRow, nullFloat64ToCSV(signal.Ping, 3))

		newRow = append(newRow, nullFloat64ToCSV(signal.SmsDeliveryTime, 3))

		newRow = append(newRow, nullInt32ToCsv(signal.RSRP))

		newRow = append(newRow, nullInt32ToCsv(signal.RSRQ))

		newRow = append(newRow, nullFloat64ToCSV(signal.Longitude, 3))

		newRow = append(newRow, nullFloat64ToCSV(signal.Latitude, 3))

		newRow = append(newRow, nullStringToCSV(signal.PCI))

		newRow = append(newRow, nullStringToCSV(signal.TAC))

		newRow = append(newRow, nullStringToCSV(signal.LAC))

		newRow = append(newRow, nullStringToCSV(signal.RAC))

		newRow = append(newRow, nullStringToCSV(signal.FrequencyBand))

		newRow = append(newRow, nullStringToCSV(signal.Arfcn))

		newRow = append(newRow, nullStringToCSV(signal.Frequency))

		newRow = append(newRow, nullStringToCSV(signal.Rscp))

		newRow = append(newRow, nullStringToCSV(signal.Ecn0))

		newRow = append(newRow, nullStringToCSV(signal.Rxlev))

		newRow = append(newRow, nullCustomTimeToCSV(signal.RecordTime))

		result = append(result, newRow)
	}

	return result
}

func nullStringToCSV(s *string) string {
	if s != nil {
		return *s
	}

	return "null"
}

func nullInt32ToCsv(i *int32) string {
	if i != nil {
		return strconv.Itoa(int(*i))
	}

	return "null"
}

func nullFloat64ToCSV(f *float64, prec int) string {
	if f != nil {
		return strconv.FormatFloat(*f, 'f', prec, 64)
	}

	return "null"
}

func nullCustomTimeToCSV(t *CustomTime) string {
	if t != nil {
		return t.Format(time.DateTime)
	}

	return "null"
}
