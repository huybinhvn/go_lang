package testing

import (
  "github.com/huybinhvn/go_lang/go-gin-restful/controller"
  "github.com/huybinhvn/go_lang/go-gin-restful/middleware"
  "github.com/gin-gonic/gin"

  "github.com/gin-contrib/sessions/cookie"
  "github.com/stretchr/testify/assert"
  "github.com/gin-contrib/sessions"
  "github.com/joho/godotenv"
  "net/http/httptest"
  "net/http"
  "testing"
  "strconv"
  "log"
  "os"
)

func loadEnv() {
  err := godotenv.Load("../.env.local")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}

func setupRouter() *gin.Engine {
  loadEnv()

  r := gin.Default()
  r.ForwardedByClientIP = true
  r.SetTrustedProxies([]string{os.Getenv("BASE_HOST")})

  sessionExp, _ := strconv.Atoi(os.Getenv("SESSION_EXP"))
  store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
  store.Options(sessions.Options{MaxAge: sessionExp})
  r.Use(sessions.Sessions(os.Getenv("SESSION_KEY"), store))
  r.GET("/incr", controller.Incr)

  r.POST("/auth/register", controller.Register)
  r.POST("/auth/login", controller.Login)

  r.Group("/api").Use(middleware.JWTAuthMiddleware())
  r.Group("/api").POST("/entry", controller.AddEntry)
  r.Group("/api").GET("/entry", controller.GetAllEntries)

  return r
}

func TestIncr(t *testing.T) {
  router := setupRouter()
  mockResponse := `{"count":0}`

  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/incr", nil)
  router.ServeHTTP(w, req)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, mockResponse, w.Body.String())
}
