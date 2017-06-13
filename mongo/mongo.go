package main //The package “main” tells the Go compiler that the package should compile as an executable program instead of a shared library.

import (
	"fmt"//implements formatted I/O with functions analogous to C's printf and scanf.
	"html/template"// implements data-driven templates for generating HTML output safe against code injection.
	// I used Go’s html/template package. The way templating works is that a html page is written with place holders 
	// for variables inserted into the html with curly brackets around them. Go code then renders the
	// .html page by filling in each of the variables as needed with the variables it has been told to use.
	"log" //implements a simple logging package. It defines a type, Logger, with methods for formatting output. 
	"net/http" //Package http provides HTTP client and server implementations.
	"gopkg.in/mgo.v2" //Package mgo offers a rich MongoDB driver for Go.
	"gopkg.in/mgo.v2/bson"
)
//Constatnts
const(
	PORT ="8082"
)
//Type to store in the db
type Person struct {
	Name  string
	Phone string
}
//Type to send to de index.html
type PageVariables struct {
	PageContacts []Person
}
//Entry point
func main() {
	//enables use of css
	fs := http.FileServer(http.Dir("static"))//returns a handler that serves HTTP requests with the contents of the file system rooted at root
    http.Handle("/static/", http.StripPrefix("/static/", fs))// you take in an http.Handler and return a new one that does something 
	//else before and/or after calling the ServeHTTP method on the original.

	http.HandleFunc("/", request)//registers the handler function for the given pattern

	fmt.Println("Running on port: ",PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
func request(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case "GET":
		getContacts(w, r)
	case "POST":
		postContacts(w, r)
	default:
		fmt.Println("Error")
	}
}
func postContacts(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("POST-INSERT")
	//gets the values to insert rom request
	r.ParseForm()
	var name = r.Form.Get("name")
	var phone = r.Form.Get("number")
	//conects to db
	session, err := mgo.Dial("localhost:27017/contacts")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	//conects to db
	c := session.DB("contacts").C("people")
	//inserts
	err = c.Insert(&Person{name, phone})
	if err != nil {
		log.Fatal(err)
	}
	getContacts(w, r)
}
func getContacts(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GET-SELECT")
	//conects to db
	session, err := mgo.Dial("localhost:27017/contacts")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//new session to connect to the table
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("contacts").C("people")
	//reads from db
	var m []bson.M
	var _ = c.Find(nil).All(&m)
	//creates slice of type Person
	slicePerson := make([]Person, 0)
	//fill the slice
	for _, v := range m {
		var name interface{} = v["name"]
		var phone interface{} = v["phone"]
		slicePerson = append(slicePerson, Person{name.(string), phone.(string)})
	}
	//data is going to send
	MyContacts := slicePerson

	MyPageVariables := PageVariables{
		PageContacts: MyContacts,
	}
	//opens index.html
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print("index parsing error: ", err)
	}
	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("index executing error: ", err)
	}
}