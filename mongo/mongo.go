package main
import (
        "html/template"
        "fmt"
	"log"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        "net/http"
)
type Person struct {
        Name string
        Phone string
}
type PageVariables struct {
        PageContacts []Person
}
func main() {
        //SERVER
	    http.HandleFunc("/", getContacts)
            //http.HandleFunc("/add", postContacts)
	    log.Fatal(http.ListenAndServe(":8082", nil))       
}
func getContacts(w http.ResponseWriter, r *http.Request){
        //MONGO
        session, err := mgo.Dial("localhost:27017/contacts")
        if err != nil {
                panic(err)
        }
        defer session.Close()

        session.SetMode(mgo.Monotonic, true)

        c := session.DB("contacts").C("people")

        var m []bson.M
        var x = c.Find(nil).All(&m)
        fmt.Println(x)
        slicePerson := make([]Person, 0)
        for i,v := range m {
                var name interface{}  = v["name"]
                var phone interface{} = v["phone"]
                slicePerson = append(slicePerson, Person{name.(string), phone.(string)})
                fmt.Println(i)
        }

        MyContacts := slicePerson

        fmt.Println("MyContacts;",slicePerson)
        MyPageVariables := PageVariables{ 
                PageContacts : MyContacts,
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
func postContacts(w http.ResponseWriter, r *http.Request){
        //MONGO
        fmt.Println("Insert data from request")

        // r.ParseForm()
        // // r.Form is now either
        // // map[animalselect:[cats]] OR
        // // map[animalselect:[dogs]]
        // // so get the animal which has been selected
        // youranimal := r.Form.Get("animalselect")

        // getContacts();

        /*session, err := mgo.Dial("localhost:27017/contacts")
        if err != nil {
                panic(err)
        }
        defer session.Close()

        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        c := session.DB("contacts").C("people")
        err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	               &Person{"Cla", "+55 53 8402 8510"})
        if err != nil {
                log.Fatal(err)
        }*/
}
