package main

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "context"
  "time"
  tea "github.com/charmbracelet/bubbletea"
  bm "github.com/charmbracelet/wish/bubbletea"
  "github.com/charmbracelet/wish"
  "github.com/gliderlabs/ssh"
)

const (
  minWidth = 80
  minHeight = 20
)

type WrongDims struct {
  width int
  height int
}

func (w WrongDims) Init() tea.Cmd {
  return nil
}

func (w WrongDims) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
      case "ctrl+c", "q":
        return w, tea.Quit
    }
  }
  return w, nil
}

func (w WrongDims) View() string {
  return fmt.Sprintf("Terminal must be %dx%d (got %dx%d)",
    minWidth, minHeight,
    w.width,
    w.height,
  )
}

type Foo struct {
  v int
}

func (f Foo) Init() tea.Cmd {
  return nil
}

func (f Foo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.String() {
        case "ctrl+c", "q":
          return f, tea.Quit
      }
    }
    return f, nil
}

func (f Foo) View() string {
  return fmt.Sprint("Oh hey: %d\n", f.v)
}

func main() {
  log.Printf("startup")
  teaOptions := []tea.ProgramOption{tea.WithAltScreen(),tea.WithOutput(os.Stderr)}
  addr := "0.0.0.0:8007"
  server, err := wish.NewServer(
    wish.WithAddress(addr),
    wish.WithIdleTimeout(10*time.Minute),
    wish.WithMiddleware(
      bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
        log.Printf("new connection (from %s) - %v", s.RemoteAddr().String(), s.Command())
        pty, _, active := s.Pty()
        if !active {
          log.Printf("no active terminal?")
          os.Exit(1)
        }
        w := pty.Window.Width
        h := pty.Window.Height
        if w < minWidth || h < minHeight {
          return WrongDims{w, h}, teaOptions
        }
        return Foo{17}, teaOptions
      }),
    ),
  )

  if err != nil {
    log.Fatalf("could not start server: %s", err.Error())
  }

  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

  log.Printf("running server: %s", server.Addr)
  go func() {
    if err := server.ListenAndServe(); err != nil {
      log.Fatalf("server returned error: %s", err.Error())
    }
  }()

  <-done
  log.Println("server done")
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer cancel()
  if err := server.Shutdown(ctx); err != nil {
    log.Fatalf("server would not shutdown gracefully: %s", err.Error())
  }
}
