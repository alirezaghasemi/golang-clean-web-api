package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type header struct {
	UserId  string
	Browser string
}

type PersonData struct {
	FirstName    string `json:"first_name" binding:"required,alpha,min=3,max=20"`
	LastName     string `json:"last_name" binding:"required,alpha,min=3,max=20"`
	MobileNumber string `json:"mobile_number" binding:"required,mobile,min=11,max=11"`
}

type TestHealthHandler struct{}

func NewTestHandler() *TestHealthHandler {
	return &TestHealthHandler{}
}

func (h *TestHealthHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "This is ok",
	})
	return
}

func (h *TestHealthHandler) HeaderBinder1(c *gin.Context) {
	userId := c.GetHeader("UserId")

	c.JSON(http.StatusOK, gin.H{
		"result": "HeaderBinder1",
		"userId": userId,
	})
}

func (h *TestHealthHandler) HeaderBinder2(c *gin.Context) {
	header := header{}
	_ = c.BindHeader(&header)

	c.JSON(http.StatusOK, gin.H{
		"result": "HeaderBinder1",
		"header": header,
	})
}

func (h *TestHealthHandler) QueryBinder1(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")

	c.JSON(http.StatusOK, gin.H{
		"result": "QueryBinder1",
		"id":     id,
		"name":   name,
	})
}

func (h *TestHealthHandler) QueryBinder2(c *gin.Context) {
	ids := c.QueryArray("id")
	name := c.Query("name")

	c.JSON(http.StatusOK, gin.H{
		"result": "QueryBinder2",
		"ids":    ids,
		"name":   name,
	})
}

func (h *TestHealthHandler) UriBinder(c *gin.Context) {
	id := c.Param("id")
	//name := c.Param("name")

	name := c.Params.ByName("name")

	c.JSON(http.StatusOK, gin.H{
		"result": "UriBinder",
		"id":     id,
		"name":   name,
	})
}

func (h *TestHealthHandler) BodyBinder(c *gin.Context) {

	p := PersonData{}

	err := c.ShouldBindJSON(&p)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"validationError": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "UriBinder",
		"person": p,
	})
}

func (h *TestHealthHandler) FormBinder(c *gin.Context) {

	p := PersonData{}

	_ = c.ShouldBind(&p)

	c.JSON(http.StatusOK, gin.H{
		"result": "UriBinder",
		"person": p,
	})
}

func (h *TestHealthHandler) FileBinder(c *gin.Context) {

	file, _ := c.FormFile("file")
	_ = c.SaveUploadedFile(file, "file")

	c.JSON(http.StatusOK, gin.H{
		"result":      "UriBinder",
		"file_name":   file.Filename,
		"file_header": file.Header,
		"file_size":   file.Size,
	})
}
