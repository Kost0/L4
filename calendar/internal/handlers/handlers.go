package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Kost0/L4/internal/business"
	"github.com/Kost0/L4/internal/models"
)

type Handler struct {
	Ch    chan *models.Event
	LogCh chan string
	Rep   *business.EventRepository
}

type result struct {
	Res []*models.Event `json:"result"`
}

func newResult(res []*models.Event, w http.ResponseWriter) {
	newRes := result{Res: res}

	buf, err := json.Marshal(&newRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		return
	}
}

type eventErr struct {
	ErrString string `json:"error"`
}

func newEventErr(err error, w http.ResponseWriter) {
	evErr := &eventErr{ErrString: err.Error()}

	buf, err := json.Marshal(evErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(evErr.ErrString, "not found") {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = w.Write(buf)
	if err != nil {
		return
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newEventErr(err, w)
		return
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			newEventErr(err, w)
			return
		}
	}()

	event, err := h.Rep.NewEvent(buf)
	if err != nil {
		newEventErr(err, w)
		return
	}

	h.Ch <- event

	w.WriteHeader(http.StatusOK)
	newResult([]*models.Event{event}, w)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newEventErr(err, w)
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	event, err := h.Rep.UpdateEvent(buf, id)
	if err != nil {
		newEventErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	newResult([]*models.Event{event}, w)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.Rep.DeleteEvent(id)
	if err != nil {
		newEventErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event deleted"))
}

func (h *Handler) GetEventForDay(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
		return
	}

	queryParams := parsedURL.Query()

	dateStr := queryParams.Get("date")
	user := queryParams.Get("user_id")

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		newEventErr(err, w)
		return
	}

	necessaryEvents, err := h.Rep.FindEventsForTime(date, user, 1)
	if err != nil {
		newEventErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	newResult(necessaryEvents, w)
}

func (h *Handler) GetEventForWeek(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
	}

	queryParams := parsedURL.Query()

	dateStr := queryParams.Get("date")
	user := queryParams.Get("user_id")

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		newEventErr(err, w)
		return
	}

	necessaryEvents, err := h.Rep.FindEventsForTime(date, user, 7)
	if err != nil {
		newEventErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	newResult(necessaryEvents, w)
}

func (h *Handler) GetEventForMonth(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		newEventErr(err, w)
	}

	queryParams := parsedURL.Query()

	dateStr := queryParams.Get("date")
	user := queryParams.Get("user_id")

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		newEventErr(err, w)
		return
	}

	necessaryEvents, err := h.Rep.FindEventsForTime(date, user, 30)
	if err != nil {
		newEventErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	newResult(necessaryEvents, w)
}
