package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	route.HandleFunc("/", homePage).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", blogPage).Methods("GET")
	route.HandleFunc("/project-detail/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/send-data-add", sendDataAdd).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", editProject).Methods("GET")
	route.HandleFunc("/send-edit-project/{index}", updateProject).Methods("POST")

	fmt.Println("Server running on port 8000")
	http.ListenAndServe("localhost:8000", route)

}

var Data = map[string]interface{}{
	"title": "Personal Web",
}

type Project struct {
	Id           int
	ProjectName  string
	StartDate    string
	EndDate      string
	Description  string
	Technologies []string
	Image        string
}

var dataProject = []Project{}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var response = map[string]interface{}{
		"Projects": dataProject,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blogPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/detail-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			BlogDetail = Project{
				ProjectName:  data.ProjectName,
				StartDate:    data.StartDate, //duration
				EndDate:      data.EndDate,
				Description:  data.Description,
				Technologies: data.Technologies,
				Image:        data.Image, //masih belum fungsi
			}
		}
	}

	data := map[string]interface{}{
		"Project": BlogDetail,
	}

	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func sendDataAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("project-name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("desc-project")
	var techno []string
	techno = r.Form["techno"]
	uploadImg := r.PostForm.Get("Imageee")

	newProject := Project{
		ProjectName:  projectName,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  description,
		Technologies: techno,
		Image:        uploadImg,
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	fmt.Println(dataProject)

	http.Redirect(w, r, "/", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf=8")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	}

	id, _ := strconv.Atoi((mux.Vars(r)["id"]))

	ProjectData := Project{}

	for index, selectProject := range dataProject {
		if id == index {
			ProjectData = Project{
				Id:           id,
				ProjectName:  selectProject.ProjectName,
				StartDate:    selectProject.StartDate,
				EndDate:      selectProject.EndDate,
				Description:  selectProject.Description,
				Technologies: selectProject.Technologies,
			}
			fmt.Println(ProjectData.Description)
		}
	}
	response := map[string]interface{}{
		"ProjectData": ProjectData,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi((mux.Vars(r)["id"]))

	projectName := r.PostForm.Get("project-name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("desc-project")
	var techno []string
	techno = r.Form["techno"]
	uploadImg := r.PostForm.Get("Imageee")

	newProject := Project{
		ProjectName:  projectName,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  description,
		Technologies: techno,
		Image:        uploadImg,
	}

	dataProject[id] = newProject

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
