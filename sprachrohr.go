// sprachrohr.go
package main


import(
    "github.com/karlyan17/jimbob"
    "github.com/gorilla/mux"
    "net/http"
    "io"
    "strconv"
    "io/ioutil"
    "html/template"
    //"sprachrohr/post"
    "sprachrohr/freshlog"
    "sprachrohr/config"
)

var CONFIG config.Config
var DB  jimbob.Bucket

func MainHandler(w http.ResponseWriter, r *http.Request) {
    nBytes,err := io.WriteString(w, "main Hellow!")
    freshlog.Debug.Print("served ", nBytes)
    if err != nil {
        freshlog.Error.Print("served %v with ", nBytes)
    }
}

func PostHandler(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    switch request.Method {
    case "GET":
        PostView(writer, request)
    case "POST":
        PostPost(writer, request)
    case "DELETE":
        PostDelete(writer, request)
    }

}

func PostView(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    if len(vars) == 0 {
        serveTemplate("posts.tmpl", writer, DB.Data)
    } else {
        id,err := strconv.Atoi(vars["id"])
        if err != nil {
            serveTemplate("posts.tmpl", writer, DB.Data)
        } else {
            serveTemplate("post.tmpl", writer, DB.Data[id])
        }
    }
}

func PostPost(writer http.ResponseWriter, request *http.Request) {
}

func PostDelete(writer http.ResponseWriter, request *http.Request) {
}

func serveTemplate(tmpl_path string, writer http.ResponseWriter, data interface{}) {
    templ_file,err := ioutil.ReadFile(tmpl_path)
    if err != nil {
        freshlog.Error.Print("failed to read template file: ", err)
    }

    templ,err := template.New(tmpl_path).Parse(string(templ_file))
    if err != nil{
        freshlog.Error.Print("failed to parse template: ", err)
        return
    }

    err = templ.Execute(writer, data)

    if err != nil {
        freshlog.Error.Print("template error ", err)
        return
    }
}

func main() {
    //
    freshlog.Debug.Print("opening jimbob bucket")
    CONFIG = config.ParseFlags()
    freshlog.Debug.Print(CONFIG)
    var err error
    DB,err = jimbob.OpenBucket(CONFIG.DB_path + "/posts")
    if err != nil {
        freshlog.Fatal.Fatal("could not open jimbob Bucket: ",err)
    }

    //for i := 0; i <10; i++ {
    //    freshlog.Debug.Print("posting to DB")
    //    _,err = DB.Post(apost)
    //    if err != nil {
    //        freshlog.Error.Print("could not POST to jimbob DB: ",err)
    //    }
    //}

    //multiplex
    r := mux.NewRouter()
    r.HandleFunc("/", MainHandler)
    r.HandleFunc("/posts", PostHandler)
    r.HandleFunc("/posts/{id:[0-9]*}", PostHandler)

    //do shit
    freshlog.Fatal.Fatal(http.ListenAndServe(CONFIG.IP + ":" + strconv.Itoa(CONFIG.Port), r))
}
