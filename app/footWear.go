package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"

	"./controller"
	"./db"

	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Server is whole server implementation for this wiki app.
// This holds database connection and router settings based on gin.
type Server struct {
	db     *sql.DB
	Engine *gin.Engine
}

// Close makes the database connection to close.
func (s *Server) Close() error {
	return s.db.Close()
}

// Init initialize server state. Connecting to database, compiling templates,
// and settings router.
func (s *Server) Init(dbconf, env string) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	db, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.db = db

	// NOTE: define helper func to use from templates here.
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"LoggedIn":    controller.LoggedIn,
		"CurrentName": controller.CurrentName,
	}).ParseGlob("templates/*"))
	s.Engine.SetHTMLTemplate(t)

	store := sessions.NewCookieStore([]byte("secretkey"))
	s.Engine.Use(sessions.Sessions("wikisession", store))
	s.Engine.Use(csrf.Middleware(csrf.Options{
		Secret: "secretkey",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismach")
			c.Abort()
		},
	}))

	s.Route()
}

// New returns server object.
func New() *Server {
	r := gin.Default()
	return &Server{Engine: r}
}

// Run starts running http server.
func (s *Server) Run(addr ...string) {
	s.Engine.Run(addr...)
}

// Route setting router for this wiki.
func (s *Server) Route() {
	item := &controller.Item{DB: s.db}
	user := &controller.User{DB: s.db}

	auth := s.Engine.Group("/")
	auth.Use(controller.AuthRequired())
	{
		auth.GET("/authtest", func(c *gin.Context) {
			c.String(200, "you're authed")
		})
		auth.GET("/new", func(c *gin.Context) {
			c.HTML(200, "new.tmpl", gin.H{
				"title":   "PictureCollection",
				"csrf":    csrf.GetToken(c),
				"context": c,
			})
		})

		auth.GET("/mypage", item.MyPageGet)
		auth.GET("/item/:id", item.Get)
		auth.GET("/item/:id/edit", item.Edit)
		auth.POST("/upload", item.Upload)
		auth.POST("/delete", item.Delete)
		auth.POST("/register", item.Make)

		auth.GET("/logout", func(c *gin.Context) {
			c.HTML(http.StatusOK, "logout.tmpl", gin.H{
				"csrf":    csrf.GetToken(c),
				"context": c,
			})
		})
		auth.POST("/logout", user.Logout)
	}

	s.Engine.GET("/", item.Root)
	s.Engine.GET("/list", item.ListGet)
	s.Engine.GET("/open/item/:id", item.OpenGet)
	s.Engine.POST("/point/:id", item.SavePoint)

	s.Engine.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"csrf": csrf.GetToken(c),
		})
	})
	s.Engine.POST("/signup", user.SignUp)

	s.Engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"csrf": csrf.GetToken(c),
		})
	})
	s.Engine.POST("/login", user.Login)
	s.Engine.Static("/images", "images")
	s.Engine.Static("/css", "css")

}

func main() {
	var (
		addr   = flag.String("addr", ":8080", "addr to bind")
		dbconf = flag.String("dbconf", "dbconfig.yml", "database configuration file.")
		env    = flag.String("env", "development", "application envirionment (production, development etc.)")
	)
	flag.Parse()
	b := New()
	b.Init(*dbconf, *env)
	b.Run(*addr)
}
