package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "crypto/md5"
    "encoding/hex"
    "html/template"
)

func main() {
    router := gin.Default()

    funcMap := make(map[string]interface{})
    funcMap["unescapeHTML"] = func(s string) template.HTML {
	    return template.HTML(s)
    }

    router.SetFuncMap(funcMap)

    router.LoadHTMLGlob("templates/*")

    router.GET("/", startPage)

    router.GET("/login", loginPage)
    router.POST("/login", loginHandler)
    router.GET("/dashboard", dashboardHandler)
    router.GET("/logout", logoutHandler)
    router.HEAD("/status")

    router.Run(":8080") // CWE-523
}

func startPage(context *gin.Context) {
    context.HTML(http.StatusOK, "index.tmpl", gin.H{
        "title": "Start page",
    })
}

func loginPage(context *gin.Context) {
    message := ""

    logout := context.Query("logout")

    if logout != "" {
        message = "<div style='border: solid 0.2em green; padding: 0.5em; margin: 0.5em 0;'> logout " + logout + "!</div>" // CWE-79: reflected via GET
    }

    context.HTML(http.StatusOK, "login.tmpl", gin.H{"message": message})
}

func dashboardHandler(context *gin.Context) {
    cookie, err := context.Cookie("user")

    if err == nil {
        context.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"user": cookie,}) // CWE-565
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

    if ("admin" == username && loginPassword == md5HashInHex(password)) {  // CWE-208
        context.SetCookie("user", username, 3600, "/", "*", false, false)
        context.Redirect(http.StatusMovedPermanently, "/dashboard")
    } else {
        context.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"message": "<div style='border: solid 0.2em red; padding: 0.5em; margin: 0.5em 0;'>Incorrect username or password for " + username + "!<div>"}) // CWE-79: reflected via POST
    }
}

func md5HashInHex(password string) string {
    hash := md5.Sum([]byte(password)) // CWE-327

    return hex.EncodeToString(hash[:])
}
