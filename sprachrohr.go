// sprachrohr.go
package main


import(
    "fmt"
    "github.com/gorilla/mux"
    "net/http"
    "io"
    "sprachrohr/post"
    "sprachrohr/freshlog"
    "sprachrohr/jimbob"
)

var db  jimbob.Bucket

func MainHandler(w http.ResponseWriter, r *http.Request) {
    nBytes,err := io.WriteString(w, "main Hellow!")
    freshlog.Debug.Print("served ", nBytes)
    if err != nil {
        freshlog.Error.Print("served %v with ", nBytes)
    }
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
    req_post := ""
    for index,rpost := range(db.Data) {
        req_post += fmt.Sprintf("%v:\tTitle: %v\tDate:%v\n%v\n", index, rpost.(post.Post)["Title"], rpost.(post.Post)["Created"], rpost.(post.Post)["Body"])
    }

    nBytes,err := io.WriteString(w, fmt.Sprintf("post Hellow!\n%v", req_post))
    freshlog.Debug.Print("served ", nBytes)
    if err != nil {
        freshlog.Error.Print("served %v with ", nBytes)
    }
}

func main() {
    //
    freshlog.Debug.Print("opening jimbob bucket")
    var err error
    db,err = jimbob.OpenBucket("db")
    if err != nil {
        freshlog.Fatal.Fatal("could not open jimbob db: ",err)
    }

    //for i := 0; i <10; i++ {
    //    freshlog.Debug.Print("posting to db")
    //    _,err = db.Post(apost)
    //    if err != nil {
    //        freshlog.Error.Print("could not POST to jimbob db: ",err)
    //    }
    //}

    //multiplex
    r := mux.NewRouter()
    r.HandleFunc("/", MainHandler)
    r.HandleFunc("/posts", PostHandler)
    r.HandleFunc("/posts/{id}", PostHandler)

    //do shit
    freshlog.Fatal.Fatal(http.ListenAndServe(":8080", r))
}
