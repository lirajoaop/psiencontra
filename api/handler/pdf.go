package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DownloadPDF(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "invalid session id")
		return
	}

	result, err := SessionSvc.GetResult(id)
	if err != nil {
		sendError(c, http.StatusNotFound, "result not found")
		return
	}

	pdfBytes, err := PDFSvc.Generate(result)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "failed to generate PDF")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=psiencontra-%s.pdf", id.String()[:8]))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
