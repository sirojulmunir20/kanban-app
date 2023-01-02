package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here
	//todo Jika data email atau password kosong
	if user.Email == "" || user.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	//todo jika err saat menggunakan userService.Login
	res, err := u.userService.Login(r.Context(), &entity.User{
		Email:    user.Email,
		Password: user.Password,
	})
	
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}

	//todo Jika user berhasil login
	responseData := map[string]interface{}{
		"user_id": res,
		"message": "login success",
	}
	
	//todo set cookei
	resStr := strconv.Itoa(res)
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   resStr,
		Path:    "/",
		Expires: time.Now().Add(5 * time.Hour),
	})
	
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseData)
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here

	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	//todo jika err saat menggunakan userService.Register
	res, err := u.userService.Register(r.Context(), &entity.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	//todo cek email sdh terdaftar
	if user.Email == res.Email && err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email already registered"))
		return
	}
	//todo jika user berhasil dibuat
	responseData := map[string]interface{}{
		"user_id": res.ID,
		"message": "register success",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responseData)

}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	//todo Menghapus cookie yang menyimpan informasi login user
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout success"})
	return
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
