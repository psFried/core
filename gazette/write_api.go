package gazette

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pippio/gazette/journal"
)

type WriteAPI struct {
	handler AppendOpHandler
}

func NewWriteAPI(handler AppendOpHandler) *WriteAPI {
	return &WriteAPI{handler: handler}
}

func (h *WriteAPI) Register(router *mux.Router) {
	router.NewRoute().Methods("PUT").HandlerFunc(h.Write)
}

func (h *WriteAPI) Write(w http.ResponseWriter, r *http.Request) {
	var op = journal.AppendOp{
		AppendArgs: journal.AppendArgs{
			Journal: journal.Name(r.URL.Path[1:]),
			Content: r.Body,
		},
		Result: make(chan journal.AppendResult, 1),
	}
	h.handler.Append(op)
	result := <-op.Result

	if result.Error != nil {
		// Return a 404 Not Found with Location header on a routing error.
		if routeError, ok := result.Error.(RouteError); ok {
			http.Redirect(w, r, routeError.RerouteURL(r.URL).String(), http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		r.Body.Close()
	} else {
		w.Header().Set(WriteHeadHeader, strconv.FormatInt(result.WriteHead, 10))
		w.WriteHeader(http.StatusNoContent)
	}
}
