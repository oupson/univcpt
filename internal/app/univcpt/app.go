package univcpt

import (
	ics "github.com/arran4/golang-ical"
	"github.com/go-chi/chi/v5"
	"github.com/oupson/univcpt/internal/pkg/calendar"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	CalendarUrl string `json:"calendar_url"`
}

type App struct {
	client *http.Client
	logger *slog.Logger

	config       *Config
	parsingRegex *regexp.Regexp

	lock      *sync.RWMutex
	calendars []*ics.Calendar
}

func (app *App) runReloadLoop() {
	for {
		app.logger.Debug("reloading calendar")
		calendars, err := calendar.GetCalendar(app.client, app.parsingRegex, app.config.CalendarUrl)
		if err != nil {
			app.logger.Warn("failed to get calendar", "error", err)
			time.Sleep(10 * time.Second)
		} else {
			app.lock.Lock()
			app.calendars = calendars
			app.lock.Unlock()

			app.logger.Info("reloaded calendar")
			time.Sleep(1 * time.Minute)
		}
	}
}

func (app *App) handleCalendar(writer http.ResponseWriter, request *http.Request) error {
	tpNbr, err := strconv.Atoi(chi.URLParam(request, "calendar"))
	if err != nil {
		return err
	}

	if tpNbr > 0 && tpNbr < 5 {
		app.lock.RLock()
		defer app.lock.RUnlock()

		if app.calendars != nil {
			writer.Header().Add("Content-Type", "text/calendar")
			writer.WriteHeader(200)
			if err := app.calendars[tpNbr-1].SerializeTo(writer); err != nil {
				return err
			}
		} else {
			writer.WriteHeader(500)
			if _, err := writer.Write([]byte("no data")); err != nil {
				return err
			}
		}
	} else {
		writer.WriteHeader(404)
		if _, err := writer.Write([]byte("No such files")); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	go app.runReloadLoop()

	router := chi.NewRouter()
	router.HandleFunc("/calendar/{calendar}", func(writer http.ResponseWriter, request *http.Request) {
		if err := app.handleCalendar(writer, request); err != nil {
			app.logger.Warn("error in request handler", "route", "calendar", "err", err)
		}
	})

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	return server.ListenAndServe()
}

func NewApp(logger *slog.Logger, icalUrl string) *App {
	client := &http.Client{}

	config := &Config{CalendarUrl: icalUrl}

	re := regexp.MustCompile(`(?m)^(Gr\s*|ANG)(TP|TD)?\s*(\d?)(FI|ALT)$`)

	return &App{
		client: client, logger: logger, config: config, parsingRegex: re, lock: new(sync.RWMutex),
	}
}
