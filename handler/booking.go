package handler

import (
	"InternBCC/database"
	"InternBCC/entity"
	"InternBCC/sdk"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Booking(c *gin.Context) {
	id, ok := c.Get("user")
	//status unautho
	if !ok {
		sdk.Fail(c, http.StatusUnauthorized, "Unauthorized")
	}
	claims := id.(entity.UserClaims)
	GeIDStr := c.Param("id")
	GeID, _ := strconv.Atoi(GeIDStr)
	//statusbad
	var req entity.BookingPost
	if err := c.BindJSON(&req); err != nil {
		sdk.FailOrError(c, http.StatusBadRequest, "Lengkapi isian Anda", err)
		return
	}
	tglStr := req.Tanggal
	template := "2006-01-02"
	parsed, err := time.Parse(template, tglStr)
	if err != nil {
		sdk.FailOrError(c, http.StatusInternalServerError, "Failed to parse", err)
		return
	}
	var get = entity.Booking{
		Nama:      req.Nama,
		Tanggal:   parsed,
		Keperluan: req.Keperluan,
		Nomer:     req.Nomor,
		Alamat:    req.Alamat,
		Fasilitas: req.Fasilitas,
		UserID:    claims.ID,
		GedungID:  uint(GeID),
	}
	req.UserID = claims.ID
	req.GedungID = uint(GeID)
	if err := database.DB.Create(&get).Error; err != nil {
		sdk.FailOrError(c, http.StatusInternalServerError, "Error to create", err)
		return
	}
	sdk.Success(c, 200, "Booking telah dibuat", req)

}

func GetBookingData(c *gin.Context) {
	type body struct {
		Nama      string `json:"nama"`
		Tanggal   string `json:"tanggal"`
		Keperluan string `json:"keperluan"`
		Nomor     string `json:"nomor"`
		Alamat    string `json:"alamat"`
	}
	id, _ := c.Get("user")
	claims := id.(entity.UserClaims)

	var get entity.Booking
	if err := database.DB.Where("user_id = ?", claims.ID).Last(&get).Error; err != nil {
		sdk.FailOrError(c, http.StatusNotFound, "Data not found", err)
		return
	}

	waktu := get.Tanggal
	format := waktu.Format("2006-01-02")

	result := body{
		Nama:      get.Nama,
		Tanggal:   format,
		Keperluan: get.Keperluan,
		Nomor:     get.Nomer,
		Alamat:    get.Alamat,
	}
	sdk.Success(c, http.StatusOK, "Data found", result)
}
