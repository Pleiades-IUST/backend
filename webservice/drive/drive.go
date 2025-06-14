package drive

import (
	"math/rand"
	"net/http"
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

	ctx.JSON(http.StatusOK, signals)
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
