package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
)

type Homework struct {
  Id      string `json:"Id"`
  Title   string `json:"Title"`
  Desc    string `json:"desc"`
  Content string `json:"content"`
}

var Homeworks []Homework

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to the HomePage!")
  fmt.Println("Endpoint Hit: homePage")
}

func returnAllHomeworks(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: returnAllHomeworks")
  json.NewEncoder(w).Encode(Homeworks)
}

func returnSingleHomework(w http.ResponseWriter, r *http.Request) {
  vars, ok := r.URL.Query()["key"]

  if !ok || len(vars[0]) < 1 {
    log.Println("Url Param 'key' is missing")
    return
  }
  key := vars[0]

  for _, homework := range Homeworks {
    if homework.Id == key {
      json.NewEncoder(w).Encode(homework)
    }
  }
}

func createNewHomework(w http.ResponseWriter, r *http.Request) {
  reqBody, _ := ioutil.ReadAll(r.Body)
  var homework Homework
  json.Unmarshal(reqBody, &homework)
  Homeworks = append(Homeworks, homework)

  json.NewEncoder(w).Encode(homework)
}

func deleteHomework(w http.ResponseWriter, r *http.Request) {
  vars, ok := r.URL.Query()["key"]

  if !ok || len(vars[0]) < 1 {
    log.Println("Url Param 'key' is missing")
    return
  }
  id := vars[0]
  for index, homework := range Homeworks {
    if homework.Id == id {
      Homeworks = append(Homeworks[:index], Homeworks[index+1:]...)
    }
  }

}

func handleRequests() {
  myRouter := http.NewServeMux()
  myRouter.HandleFunc("/", homePage)
  myRouter.HandleFunc("/homeworks", returnAllHomeworks)
  myRouter.HandleFunc("/homework", createNewHomework)
  myRouter.HandleFunc("/homework/{id}", returnSingleHomework)
  log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func main() {
  Homeworks = []Homework{
    Homework{Id: "1", Title: "Hello", Desc: "homework Description", Content: "homework Content"},
    Homework{Id: "2", Title: "Hello 2", Desc: "homework Description", Content: "homework Content"},
  }
  handleRequests()
}