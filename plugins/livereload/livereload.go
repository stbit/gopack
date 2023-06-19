package livereload

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/logger"
	"github.com/stbit/gopack/plugins"
)

var pluginName = "livereload"

type Options struct {
	Address string
}

func New(opt Options) plugins.PluginRegister {
	return func(m *plugins.ManagerContext) error {
		if !m.IsWatch() {
			return nil
		}

		updatedAt := time.Now()
		hub := newHub()
		go hub.run()

		sendChange := func() {
			hub.broadcast <- []byte("changed")
		}

		go func() {
			if opt.Address == "" {
				opt.Address = "localhost:32328"
			}

			http.HandleFunc("/livereload.js", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/javascript")
				w.Write([]byte(fmt.Sprintf(`
					const start = function() {
						var conn = new WebSocket("ws://%s/ws");

						conn.addEventListener('close', function() {
							setTimeout(start, 1000);
						});
						conn.onmessage = function () {
							window.location.reload()
						};
					}

					start();
				`, opt.Address)))
			})

			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				serveWs(hub, w, r, updatedAt)
			})

			log.Printf("%s %s", logger.Magenta("http livereload script:"), logger.Green(opt.Address+"/livereload.js"))

			if err := http.ListenAndServe(opt.Address, nil); err != nil {
				log.Fatal(err)
			}
		}()

		m.AddHook(pluginName, hooks.HOOK_AFTER_COMPILE, func() error {
			sendChange()
			updatedAt = time.Now()
			return nil
		})

		return nil
	}
}
