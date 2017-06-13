package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

type Person struct {
	Name  string
	Phone string
}
type PageVariables struct {
	PageContacts []Person
}

func main() {
	//SERVER
	fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", request)
	//     http.HandleFunc("/add", postContacts)
	//     http.HandleFunc("/delete", deleteContact)
	fmt.Println("Running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
func request(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//fmt.Println(r)
	switch r.Method {
	case "GET":
		getContacts(w, r)
	// Serve the resource.
	case "POST":
		postContacts(w, r)
	// Create a new record.
	case "PUT":
		updateContact(w, r)
	// Update an existing record.
	case "DELETE":
		deleteContact(w, r)
	// Remove the record.
	default:
		fmt.Println("Error")
	}
}
func postContacts(w http.ResponseWriter, r *http.Request) {
	//MONGO
	fmt.Println("POST-INSERT")
	r.ParseForm()
	fmt.Println(r)
	var name = r.Form.Get("name")
	var phone = r.Form.Get("number")
	// fmt.Println(name)
	// fmt.Println(phone)

	session, err := mgo.Dial("localhost:27017/contacts")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("contacts").C("people")
	err = c.Insert(&Person{name, phone})
	if err != nil {
		log.Fatal(err)
	}
	getContacts(w, r)
}
func getContacts(w http.ResponseWriter, r *http.Request) {
	//MONGO
	fmt.Println("GET-SELECT")
	session, err := mgo.Dial("localhost:27017/contacts")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("contacts").C("people")

	var m []bson.M
	var _ = c.Find(nil).All(&m)
	// fmt.Println(x)
	slicePerson := make([]Person, 0)
	for _, v := range m {
		var name interface{} = v["name"]
		var phone interface{} = v["phone"]
		slicePerson = append(slicePerson, Person{name.(string), phone.(string)})
	}

	MyContacts := slicePerson

	// fmt.Println("MyContacts;",slicePerson)
	MyPageVariables := PageVariables{
		PageContacts: MyContacts,
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print("index parsing error: ", err)
	}
	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("index executing error: ", err)
	}
}
func updateContact(w http.ResponseWriter, r *http.Request) {
	//MONGO
	fmt.Println("PUT-Update")
	fmt.Println(r)
	todo := &Person{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(todo)
	if err != nil {
		fmt.Println("ERROR decoding JSON - ", err)
		return
	}
	fmt.Println(todo)
}
func deleteContact(w http.ResponseWriter, r *http.Request) {
	//MONGO
	fmt.Println("DELETE-Delete")
	fmt.Println(r)
	todo := &Person{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(todo)
	if err != nil {
		fmt.Println("ERROR decoding JSON - ", err)
		return
	}
	fmt.Println(todo)
}
