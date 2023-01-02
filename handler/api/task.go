package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	usrId := r.Context().Value("id").(string)

	if usrId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	usrIdInt, _ := strconv.Atoi(usrId)
	taskId := r.URL.Query().Get("task_id")
	if taskId == "" {
		res, err := t.taskService.GetTasks(r.Context(), usrIdInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	} else {
		taskIdInt, _ := strconv.Atoi(taskId)
		res, err := t.taskService.GetTaskByID(r.Context(), taskIdInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	// TODO: answer here
	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	usrId := r.Context().Value("id").(string)
	if usrId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	usrIdInt, _ := strconv.Atoi(usrId)
	res, err := t.taskService.StoreTask(r.Context(), &entity.Task{
		UserID:      usrIdInt,
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	responseData := map[string]interface{}{
		"user_id": res.UserID,
		"task_id": res.ID,
		"message": "success create new task",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responseData)
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	usrId := r.Context().Value("id").(string)
	taskID := r.URL.Query().Get("task_id")
	taskIdInt, _ := strconv.Atoi(taskID)
	usrIdInt, _ := strconv.Atoi(usrId)

	err := t.taskService.DeleteTask(r.Context(), taskIdInt)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	responseData := map[string]interface{}{
		"user_id": usrIdInt,
		"task_id": taskIdInt,
		"message": "success delete task",
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseData)

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here
	usrId := r.Context().Value("id").(string)
	if usrId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	usrIdInt, _ := strconv.Atoi(usrId)
	res, err := t.taskService.UpdateTask(r.Context(), &entity.Task{
		UserID:      usrIdInt,
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	responseData := map[string]interface{}{
		"user_id": res.UserID,
		"task_id": res.ID,
		"message": "success update task",
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseData)

}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
