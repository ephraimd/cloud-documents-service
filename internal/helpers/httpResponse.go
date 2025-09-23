package helpers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func RespondOk(message string, data_ *gin.H) (int, gin.H) {
	responseTemplate := gin.H{
		"data":    data_,
		"message": message,
	}
	return 200, responseTemplate
}

func RespondCreated(message string, data_ *gin.H) (int, gin.H) {
	responseTemplate := gin.H{
		"data":    data_,
		"message": message,
	}
	return 201, responseTemplate
}

func RespondError(message string, data_ *gin.H, code int) (int, gin.H) {
	var responseTemplate gin.H
	if data_ != nil {
		responseTemplate = gin.H{
			"data":    data_,
			"message": message,
		}
	} else {
		responseTemplate = gin.H{
			"message": message,
		}
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	return code, responseTemplate
}
