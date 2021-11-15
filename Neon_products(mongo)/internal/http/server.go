package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"neon_products/internal/models"
	"neon_products/internal/store"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/neons", func(w http.ResponseWriter, r *http.Request) {
		neon := new(models.Neon)
		if err := json.NewDecoder(r.Body).Decode(neon); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		neon.ID = primitive.NewObjectID()
		s.store.Create(r.Context(), neon)
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/neons", func(w http.ResponseWriter, r *http.Request) {
		neons, err := s.store.All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		render.JSON(w, r, neons)
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/neons/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		neon, err := s.store.ByID(r.Context(), idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		render.JSON(w, r, neon)
		w.WriteHeader(http.StatusOK)
	})

	r.Put("/neons", func(w http.ResponseWriter, r *http.Request) {
		neon := new(models.Neon)
		if err := json.NewDecoder(r.Body).Decode(&neon); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		err := s.store.Update(r.Context(), neon)
		if err != nil {
			fmt.Fprintf(w, "err: %v", err)
			w.WriteHeader(http.StatusNotModified)
		}
		w.WriteHeader(http.StatusOK)

	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		if err := s.store.Delete(r.Context(), idStr); err != nil {
			fmt.Fprintf(w, "err: %v", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
