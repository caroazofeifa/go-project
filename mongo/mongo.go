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
         //MONGO
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
        }
        */
        //SERVER
	    http.HandleFunc("/", index)
	    log.Fatal(http.ListenAndServe(":8082", nil))       
}
func index(w http.ResponseWriter, r *http.Request){
        //MONGO
        fmt.Println("Empieza el index")
        session, err := mgo.Dial("localhost:27017/contacts")
        if err != nil {
                panic(err)
        }
        defer session.Close()

        session.SetMode(mgo.Monotonic, true)

        c := session.DB("contacts").C("people")

        // result := Person{}
        // err = c.Find(nil)
        // if err != nil {
        //         log.Fatal(err)
        // }

        var m []bson.M
        var x = c.Find(nil).All(&m)
        fmt.Println(x)
        
        MyContacts:=[]Person{
            for i, v := range m {
                Person{v["name"], "1234"}
                //fmt.Println(v["name"])
                fmt.Println(i)
            }
                //Person{"Alicia", "1234"},
                //Person{"Aleman", "5678"},   
        }
        MyPageVariables := PageVariables{ //store the date and time in a struct
                PageContacts : MyContacts,
        }

        t, err := template.ParseFiles("index.html") //parse the html file homepage.html
        if err != nil { // if there is an error
                log.Print("index parsing error: ", err) // log it
                }
        err = t.Execute(w, MyPageVariables) //execute the index and pass it the HomePageVars struct to fill in the gaps
        //err = t.Execute(w, MyContacts) //execute the index and pass it the HomePageVars struct to fill in the gaps
        if err != nil { // if there is an error
                log.Print("index executing error: ", err) //log it
  	}
}
