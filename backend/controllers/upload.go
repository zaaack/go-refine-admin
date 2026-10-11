package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Upload saves a posted file and records it (regular mode, file business).
//
//	@Summary	Upload a file
//	@Tags		Upload
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		file	formData	file	true	"file"
//	@Success	200		{object}	response.Envelope
//	@Router		/upload [post]
//	@Security	Bearer
func Upload(gdb *gorm.DB, dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		f, err := c.FormFile("file")
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		dst := filepath.Join(dir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(f.Filename)))
		if err := c.SaveUploadedFile(f, dst); err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		rec := models.Upload{
			Name: f.Filename,
			Url:  "/" + filepath.ToSlash(dst),
			Size: f.Size,
			Ext:  filepath.Ext(f.Filename),
		}
		rec.SetCreateBy(auth.GetUserID(c))
		if err := gdb.Create(&rec).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, rec, "")
	}
}
