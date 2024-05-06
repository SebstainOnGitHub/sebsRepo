package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var pl = fmt.Println

type ToDoList struct {
	ToDoCount int
	ToDos     []string
}

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func write(writer http.ResponseWriter, msg string) {
	_, err := writer.Write([]byte(msg))
	errorCheck(err)
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	errorCheck(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	errorCheck(scanner.Err())
	return lines
}

func englishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Hello Internet!")
}

func germanHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Hallo Internet!")
}

func interactHandler(writer http.ResponseWriter, request *http.Request) {
	toDoVals := getStrings("todos.txt")

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))
	errorCheck(err)
	fmt.Printf("%#v\n", toDoVals)
	tmpl, err := template.ParseFiles("view.html")
	errorCheck(err)

	todos := ToDoList{
		ToDoCount: len(toDoVals),
		ToDos:     toDoVals,
	}

	defer file.Close()

	err = tmpl.Execute(writer, todos)
	errorCheck(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("new.html")
	errorCheck(err)
	err = tmpl.Execute(writer, nil)
	errorCheck(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	todo := request.FormValue("todo")
	if todo == "" {
		http.Redirect(writer, request, "/interact", http.StatusFound)
	}

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))
	errorCheck(err)

	_, err = fmt.Fprintln(file, todo)

	errorCheck(err)

	defer file.Close()

	http.Redirect(writer, request, "/interact", http.StatusFound)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	toDoVals := getStrings("todos.txt")
	
	tmpl, err := template.ParseFiles("delete.html")
	errorCheck(err)

	if len(toDoVals) == 0 {
		http.Redirect(writer, request, "/interact", http.StatusFound)
	}

	fmt.Printf("%#v\n", toDoVals)
	errorCheck(err)

	todos := ToDoList{
		ToDoCount: len(toDoVals),
		ToDos:     toDoVals,
	}

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))

	errorCheck(err)

	os.Truncate("todos.txt", 0)

	defer file.Close()

	for i := range len(toDoVals) - 1{
		fmt.Fprintf(file, toDoVals[i] + "\n")
	}

	err = tmpl.Execute(writer, todos)

	errorCheck(err)
}


func main() {
	http.HandleFunc("/hello", englishHandler)
	http.HandleFunc("/hallo", germanHandler)
	http.HandleFunc("/interact", interactHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
