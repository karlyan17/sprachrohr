// sprachrohr.go
package main


import(
    "fmt"
    "github.com/gorilla/mux"
    "net/http"
    "io"
    "time"
    "sprachrohr/post"
    "sprachrohr/freshlog"
    "sprachrohr/jimbob"
)

var Posts map[int]post.Post

var apost = post.Post{
        time.Now(),
        "this is title",
        "this is body",
        [5]string{"first","second"},
    }

func MainHandler(w http.ResponseWriter, r *http.Request) {
    nBytes,err := io.WriteString(w, "main Hellow!")
    freshlog.Debug.Print("served ", nBytes)
    if err != nil {
        freshlog.Error.Print("served %v with ", nBytes)
    }
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
    nBytes,err := io.WriteString(w, fmt.Sprintf("post Hellow!\n%v", apost))
    freshlog.Debug.Print("served ", nBytes)
    if err != nil {
        freshlog.Error.Print("served %v with ", nBytes)
    }
}

func main() {
    //
    Posts = make(map[int]post.Post)
    Posts[1] = apost
    db := jimbob.NewBucket("db.jb", Posts)
    freshlog.Debug.Println(db)

    //multiplex
    r := mux.NewRouter()
    r.HandleFunc("/", MainHandler)
    r.HandleFunc("/posts", PostHandler)
    r.HandleFunc("/posts/{id}", PostHandler)

    //do shit
    freshlog.Fatal.Fatal(http.ListenAndServe(":8080", r))
}
