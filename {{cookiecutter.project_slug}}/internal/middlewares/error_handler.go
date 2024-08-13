package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/utils/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var logger = configs.GetLogger()

var (
	ErrorInternalError = errors.New("whoops! something went wrong")
)

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = UcFirst(word)
		} else {
			entries[index] = LcFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}

func ValidationErrorToText(e validator.FieldError) string {
	word := Split(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param())
	}
	return fmt.Sprintf("%s is not valid", word)
}

// This method collects all errors and submits them to Rollbar
func ErrorHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {

				logger.Errorf(e.Error())
				switch e.Type {
				case gin.ErrorTypePublic:
					var apiError *responses.APIError
					if errors.As(e.Err, &apiError) {
						if !c.Writer.Written() {
							c.JSON(apiError.Status, gin.H{"code": apiError.ErrorCode, "message": apiError.Message})
						}
					}
				case gin.ErrorTypeBind:
					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}

					var errs validator.ValidationErrors
					if errors.As(e.Err, &errs) {
						list := make(map[string]string)
						for _, err := range errs {
							list[strings.ToLower(err.Field())] = ValidationErrorToText(err)
						}

						c.JSON(status, gin.H{"code": "ErrBadRequest", "message": "Invalid request", "details": list})
					} else {
						c.JSON(status, gin.H{"code": "ErrBadRequest", "message": "Invalid request"})
					}

					// default:
					// 	logger.Errorf("Unexpected error %s", e.Err)
				}

			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"code": "ErrUnexpected", "message": ErrorInternalError.Error()})
			}
		}
	}
}
