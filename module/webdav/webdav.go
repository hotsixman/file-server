package webdav

import (
	"file-server/module/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type WebdavApp struct {
	engine      *gin.Engine
	dav         *webdav.Handler
	authManager *auth.Auth
}

func NewWebdavApp(rootDir string, lockSystem webdav.LockSystem, authManager *auth.Auth) *WebdavApp {
	app := &WebdavApp{
		engine: gin.Default(),
		dav: &webdav.Handler{
			Prefix:     "/dav",
			FileSystem: webdav.Dir(rootDir),
			LockSystem: lockSystem,
		},
		authManager: authManager,
	}

	app.registerHandlers()

	return app
}

func (app *WebdavApp) Listen(address string) error {
	return app.engine.Run(address)
}

func (app *WebdavApp) registerHandlers() {
	methods := []string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS",
		"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	}
	for _, method := range methods {
		app.engine.Handle(method, "/dav", func(c *gin.Context) {
			user, password, ok := c.Request.BasicAuth()
			if !ok {
				user = c.Request.URL.User.Username()
				password, _ = c.Request.URL.User.Password()
			}

			if !app.checkPermission(user, password, c.Request.URL.Path) {
				c.Status(403)
				return
			}

			app.dav.ServeHTTP(c.Writer, c.Request)
		})
		app.engine.Handle(method, "/dav/*path", func(c *gin.Context) {
			user, password, ok := c.Request.BasicAuth()
			if !ok {
				user = c.Request.URL.User.Username()
				password, _ = c.Request.URL.User.Password()
			}

			if !app.checkPermission(user, password, c.Request.URL.Path) {
				c.Status(403)
				return
			}

			app.dav.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func (app *WebdavApp) checkPermission(user, password, path string) bool {
	if app.authManager == nil {
		return true
	}

	if !app.authManager.Authenticate(user, password) {
		return false
	}

	path, _ = strings.CutPrefix(path, "/dav")
	return app.authManager.CheckPermission(user, path)
}
