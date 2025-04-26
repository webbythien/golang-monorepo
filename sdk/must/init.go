package must

import "github.com/webbythien/monorepo/pkg/l"

var (
	ll = l.New()
)

func init() {
	ll.Info("Initializing must package")
}
