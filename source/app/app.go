package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func loadTemplates() multitemplate.Renderer {

	funcMap := make(map[string]interface{})
	funcMap["unescapeHTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	renderer := multitemplate.NewRenderer()
	renderer.AddFromFiles("index", "templates/layout.html", "templates/index.html")
	renderer.AddFromFilesFuncs("login", funcMap, "templates/layout.html", "templates/login.html")
	renderer.AddFromFilesFuncs("dashboard", funcMap, "templates/layout.html", "templates/dashboard.html")

	return renderer
}

func main() {
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "favicon.ico")
	router.HTMLRender = loadTemplates()

	router.GET("/", startPage)

	router.GET("/login", loginPage)
	router.POST("/login", loginHandler)
	router.GET("/dashboard", dashboardHandler)
	router.GET("/logout", logoutHandler)
	router.HEAD("/status")

	router.Run(":8080") // CWE-523
}

func startPage(context *gin.Context) {
	context.HTML(http.StatusOK, "index", gin.H{
		"title": "Start page",
	})
}

func loginPage(context *gin.Context) {
	message := ""

	logout := context.Query("logout")

	if logout != "" {
		message = "<div class='success'>Logged out " + logout + "!</div>" // CWE-79: reflected via GET
	}

	context.HTML(http.StatusOK, "login", gin.H{"message": message, "title": "Login"})
}

func dashboardHandler(context *gin.Context) {
	cookie, err := context.Cookie("user")

	if err == nil {
		context.HTML(http.StatusOK, "dashboard", gin.H{"user": cookie, "title": "Dashboard"}) // CWE-565
	} else {
		context.Redirect(http.StatusMovedPermanently, "/login")
	}

}

func logoutHandler(context *gin.Context) {
	cookie, err := context.Cookie("user")

	if err == nil {
		// delete cookie
		context.SetCookie("user", cookie, -1, "/", "*", false, false)
	}

	context.Redirect(http.StatusMovedPermanently, "/login?logout=successful")
}

func loginHandler(context *gin.Context) {
	loginPassword := "21232f297a57a5a743894a0e4a801fc3" // CWE-257

	username := context.PostForm("username")
	password := context.PostForm("password")

	if "admin" == username && loginPassword == md5HashInHex(password) { // CWE-208
		context.SetCookie("user", username, 3600, "/", "*", false, false)
		context.Redirect(http.StatusMovedPermanently, "/dashboard")
	} else {
		context.HTML(http.StatusUnauthorized, "login", gin.H{"message": "<div class='error'>Incorrect username or password for " + username + "!<div>", "title": "Login"}) // CWE-79: reflected via POST
	}
}

func md5HashInHex(password string) string {
	hash := md5.Sum([]byte(password)) // CWE-327

	return hex.EncodeToString(hash[:])
}
