package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"net/http"
	"refactoring/pkg/models"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const users = `users.json`

type (
	userRequest struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}
)

var (
	mutex sync.RWMutex
)

func (c *userRequest) Bind(r *http.Request) error { return nil }

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().String()))
}

func (app *application) showUsers(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	file, err := ioutil.ReadFile(users)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userStore := models.UserStore{}
	err = json.Unmarshal(file, &userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	render.JSON(w, r, userStore.List)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	file, err := ioutil.ReadFile(users)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userStore := models.UserStore{}
	err = json.Unmarshal(file, &userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := chi.URLParam(r, "id")
	if _, ok := userStore.List[id]; !ok {
		err = render.Render(w, r, badRequestError(ErrUserNotFound))
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	user := userStore.List[id]
	render.JSON(w, r, user)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := ioutil.ReadFile(users)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userStore := models.UserStore{}
	err = json.Unmarshal(file, &userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	request := userRequest{}

	if err := render.Bind(r, &request); err != nil {
		err = render.Render(w, r, badRequestError(err))
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	userStore.UserId++
	user := models.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.Email,
	}

	id := strconv.Itoa(userStore.UserId)
	userStore.List[id] = user

	json, err := json.Marshal(&userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ioutil.WriteFile(users, json, fs.ModePerm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{
		"user_id": id,
	})
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := ioutil.ReadFile(users)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userStore := models.UserStore{}
	err = json.Unmarshal(file, &userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	request := userRequest{}

	if err := render.Bind(r, &request); err != nil {
		err = render.Render(w, r, badRequestError(err))
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := userStore.List[id]; !ok {
		err = render.Render(w, r, badRequestError(ErrUserNotFound))
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	user := userStore.List[id]
	user.DisplayName = request.DisplayName
	userStore.List[id] = user

	json, err := json.Marshal(&userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ioutil.WriteFile(users, json, fs.ModePerm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := ioutil.ReadFile(users)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userStore := models.UserStore{}
	err = json.Unmarshal(file, &userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := userStore.List[id]; !ok {
		err = render.Render(w, r, badRequestError(ErrUserNotFound))
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	delete(userStore.List, id)

	json, err := json.Marshal(&userStore)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ioutil.WriteFile(users, json, fs.ModePerm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
