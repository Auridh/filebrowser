package auth

import (
	"errors"
	fbErrors "github.com/Auridh/filebrowser/v2/errors"
	"github.com/Auridh/filebrowser/v2/settings"
	"github.com/Auridh/filebrowser/v2/users"
	"github.com/sethvargo/go-password/password"
	"net/http"
)

// MethodProxyAuth is used to identify no auth.
const MethodProxyAuth settings.AuthMethod = "proxy"

// ProxyAuth is a proxy implementation of an auther.
type ProxyAuth struct {
	Header string `json:"header"`
}

// Auth authenticates the user via an HTTP header.
func (a ProxyAuth) Auth(r *http.Request, usr users.Store, set *settings.Settings, srv *settings.Server) (*users.User, error) {
	username := r.Header.Get(a.Header)
	user, err := usr.Get(srv.Root, username)
	pass, _ := password.Generate(64, 10, 10, false, true)

	if errors.Is(err, fbErrors.ErrNotExist) {
		err := usr.Save(&users.User{
			Username:     username,
			Password:     pass,
			Scope:        set.Defaults.Scope,
			Locale:       set.Defaults.Locale,
			ViewMode:     set.Defaults.ViewMode,
			SingleClick:  set.Defaults.SingleClick,
			Sorting:      set.Defaults.Sorting,
			Perm:         set.Defaults.Perm,
			Commands:     set.Defaults.Commands,
			HideDotfiles: set.Defaults.HideDotfiles,
		})
		if err != nil {
			return nil, err
		}

		user, err = usr.Get(srv.Root, username)
	}

	return user, err
}

// LoginPage tells that proxy auth doesn't require a login page.
func (a ProxyAuth) LoginPage() bool {
	return false
}
