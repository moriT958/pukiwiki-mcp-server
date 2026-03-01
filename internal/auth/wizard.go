package auth

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

//go:embed templates
var templatesFS embed.FS

// 初期設定に使うデータ
type setupData struct {
	URL      string
	Username string
	Scope    string
	Error    string
}

// 初期設定用の WebUI を起動する
func runcWizard(ctx context.Context) (*config, error) {
	setupTmpl, doneTmpl, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	ln, err := listenLocal()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{}) // /done へのアクセス時に close される
	result := make(chan *config, 1)
	mux := newServeMux(setupTmpl, doneTmpl, done, result)

	srv := &http.Server{Handler: mux}
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	addr := ln.Addr().String()
	browserURL := "http://" + addr
	fmt.Fprintf(os.Stderr, "pukiwiki-mcp: setup wizard at %s\n", browserURL)

	if err := exec.Command("open", browserURL).Start(); err != nil {
		fmt.Fprintf(os.Stderr, "pukiwiki-mcp: could not open browser. Please visit: %s\n", browserURL)
	}

	select {
	case <-done:
	case serveErr := <-errCh:
		return nil, fmt.Errorf("setup wizard server error: %w", serveErr)
	case <-ctx.Done():
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "pukiwiki-mcp: setup wizard shutdown error: %v\n", err)
		}
		return nil, ctx.Err()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "pukiwiki-mcp: setup wizard shutdown error: %v\n", err)
	}

	return <-result, nil
}

// 設定ウィザードを開くポート
const wizardPort = "8742"

func parseTemplates() (*template.Template, *template.Template, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/setup.html", "templates/done.html")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return tmpl.Lookup("setup.html"), tmpl.Lookup("done.html"), nil
}

func listenLocal() (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", wizardPort))
	if err != nil {
		// OS 割り当てにフォールバック
		ln, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("cannot open local port for setup wizard: %w", err)
		}
	}
	return ln, nil
}

func newServeMux(setupTmpl, doneTmpl *template.Template, done chan struct{}, result chan<- *config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", setupPageHandler(setupTmpl))
	mux.HandleFunc("POST /submit", submitHandler(setupTmpl, result))
	mux.HandleFunc("GET /done", doneHandler(doneTmpl, done))
	return mux
}

func setupPageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, setupData{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func submitHandler(setupTmpl *template.Template, result chan<- *config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		rawURL := r.FormValue("url")
		username := r.FormValue("username")
		password := r.FormValue("password")
		scope := r.FormValue("scope")

		renderErr := func(msg string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := setupTmpl.Execute(w, setupData{URL: rawURL, Username: username, Scope: scope, Error: msg}); err != nil {
				fmt.Fprintf(os.Stderr, "pukiwiki-mcp: failed to render setup template: %v\n", err)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}

		if rawURL == "" || username == "" {
			renderErr("URL とユーザー名は必須です。")
			return
		}

		cfg := &config{
			URL:      rawURL,
			Username: username,
			Password: password,
			Scope:    scope,
		}
		if err := Save(cfg); err != nil {
			renderErr(fmt.Sprintf("保存に失敗しました: %v", err))
			return
		}

		result <- cfg
		http.Redirect(w, r, "/done", http.StatusSeeOther)
	}
}

func doneHandler(doneTmpl *template.Template, done chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := doneTmpl.Execute(w, nil); err != nil {
			fmt.Fprintf(os.Stderr, "pukiwiki-mcp: failed to render done template: %v\n", err)
			fmt.Fprintln(w, "Setup complete. You can close this tab.")
		}

		go func() {
			select {
			case <-done:
			default:
				close(done)
			}
		}()
	}
}
