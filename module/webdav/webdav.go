package webdav

import (
	fs "file-server/module/Fs"
	"file-server/module/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type WebdavApp struct {
	engine      *gin.Engine
	authManager auth.AuthManager
	lockSystem  webdav.LockSystem
	davMap      map[string]*webdav.Handler
	rootDir     string
}

func NewWebdavApp(rootDir string, lockSystem webdav.LockSystem, authManager auth.AuthManager) *WebdavApp {
	app := &WebdavApp{
		engine:      gin.Default(),
		authManager: authManager,
		lockSystem:  lockSystem,
		davMap:      map[string]*webdav.Handler{},
		rootDir:     rootDir,
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
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				username = c.Request.URL.User.Username()
				password, _ = c.Request.URL.User.Password()
			}

			if !app.Permission(username, password, c.Request.URL.Path) {
				c.Status(403)
				return
			}

			app.getWebdavHandler(username).ServeHTTP(c.Writer, c.Request)
		})
		app.engine.Handle(method, "/dav/*path", func(c *gin.Context) {
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				username = c.Request.URL.User.Username()
				password, _ = c.Request.URL.User.Password()
			}

			if !app.Permission(username, password, c.Request.URL.Path) {
				c.Status(403)
				return
			}

			app.getWebdavHandler(username).ServeHTTP(c.Writer, c.Request)
		})
	}
}

func (app *WebdavApp) Permission(user, password, path string) bool {
	if app.authManager == nil {
		return true
	}

	if !app.authManager.Authenticate(user, password) {
		return false
	}

	path, _ = strings.CutPrefix(path, "/dav")
	return app.authManager.Permission(user, path)
}

func (app *WebdavApp) getWebdavHandler(username string) *webdav.Handler {
	dav, ok := app.davMap[username]
	if !ok {
		dav = &webdav.Handler{
			Prefix:     "/dav",
			FileSystem: fs.NewFs(app.authManager.DirMap(app.rootDir, username), app.authManager, username).ToWebdavFs(),
			LockSystem: app.lockSystem,
		}
		app.davMap[username] = dav
	}
	return dav
}
