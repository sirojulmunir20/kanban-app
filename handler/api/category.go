package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	id := r.Context().Value("id").(string)

	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	idInt, _ := strconv.Atoi(id)
	categories, err := c.categoryService.GetCategories(r.Context(), idInt)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(categories)

}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	// TODO: answer here

	if category.Type == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	id := r.Context().Value("id").(string)
	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	usrId,_:= strconv.Atoi(id)
	res, err := c.categoryService.StoreCategory(r.Context(), &entity.Category{
		Type: category.Type,
		UserID: usrId,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	responseData := map[string]interface{}{
		"user_id":     res.UserID,
		"category_id": res.ID,
		"message":     "success create new category",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responseData)

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	usrId := r.Context().Value("id").(string)
	if usrId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	catId := r.URL.Query().Get("category_id")

	if catId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteCatId, _ := strconv.Atoi(catId) 
	err := c.categoryService.DeleteCategory(r.Context(),deleteCatId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	usrIdInt,_:= strconv.Atoi(usrId)
	responseData := map[string]interface{}{
		"user_id": usrIdInt,
		"category_id" : deleteCatId,
		"message": "success delete category",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
	
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
