package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teuber789/chore-tracker/internal/gen"
)

func NewHttpRouter(store ChoreTrackerStore) (*mux.Router, error) {
	s := &httpSrv{store: store}

	r := mux.NewRouter()
	r.HandleFunc("/families", s.addFamily).Methods("POST")
	r.HandleFunc("/children", s.addChild).Methods("POST")
	r.HandleFunc("/chores", s.createChore).Methods("POST")
	r.HandleFunc("/chores/{id}", s.deleteChore).Methods("DELETE")
	r.HandleFunc("/chores", s.getChores).Methods("GET")
	r.HandleFunc("/completions", s.getCompletedChores).Methods("GET")
	r.HandleFunc("/completions/mark", s.markChoreCompleted).Methods("POST")

	return r, nil
}

type httpSrv struct {
	store ChoreTrackerStore
}

// IRL, we would serialize and deserialize to an intermediate representation instead of GRPC structs
func (s *httpSrv) addFamily(w http.ResponseWriter, r *http.Request) {
	var body gen.AddFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fam, err := s.store.AddFamily(context.TODO(), &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&fam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpSrv) addChild(w http.ResponseWriter, r *http.Request) {
	var body gen.AddChildRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	child, err := s.store.AddChild(context.TODO(), &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&child)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpSrv) createChore(w http.ResponseWriter, r *http.Request) {
	var body gen.CreateChoreRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chore, err := s.store.CreateChore(context.TODO(), &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&chore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// IRL, this would have permissions checks to ensure someone doesn't delete
// a chore they aren't supposed to.
func (s *httpSrv) deleteChore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.store.DeleteChore(context.TODO(), id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *httpSrv) getChoresRequest(r *http.Request) (*gen.GetChoresRequest, error) {
	q := r.URL.Query()

	// IRL, handle failures better
	token := q.Get("pageToken")
	sizeStr := q.Get("pageSize")
	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, err
	}

	size := uint32(sizeInt)
	p := gen.Pageable{PageToken: token, PageSize: size}

	famId, err := strconv.ParseUint(q.Get("familyId"), 10, 64)
	if err != nil {
		return nil, err
	}

	childId, err := strconv.ParseUint(q.Get("childId"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &gen.GetChoresRequest{Pageable: &p, FamilyId: famId, ChildId: childId}, nil
}

// IRL, this handler would do more error checking
func (s *httpSrv) getChores(w http.ResponseWriter, r *http.Request) {
	req, err := s.getChoresRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.store.GetChores(context.TODO(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpSrv) getCompletedChores(w http.ResponseWriter, r *http.Request) {
	req, err := s.getChoresRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.store.GetCompletedChores(context.TODO(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpSrv) markChoreCompleted(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// IRL, handle failures better
	famId, err := strconv.ParseUint(q.Get("familyId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	childId, err := strconv.ParseUint(q.Get("childId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	choreId, err := strconv.ParseUint(q.Get("choreId"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := gen.MarkChoreCompletedRequest{FamilyId: famId, ChildId: childId, ChoreId: choreId}
	err = s.store.MarkChoreCompleted(context.TODO(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
