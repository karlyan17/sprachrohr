// sprachrohr.go
package main


import(
    "github.com/karlyan17/jimbob"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "io/ioutil"
    "html/template"
    //"sprachrohr/post"
    "sprachrohr/freshlog"
    "sprachrohr/config"
)

var CONFIG config.Config
var DB  jimbob.Bucket

func MainHandler(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    http.Redirect(writer, request, "/posts", http.StatusSeeOther)
}

func PostsViewer(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    serveTemplate("posts.tmpl", writer, DB.Data)
}

func PostViewer(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    if len(vars) == 0 {
        http.Redirect(writer, request, "/posts", http.StatusSeeOther)
        return
    }

    id,err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Redirect(writer, request, "/posts", http.StatusSeeOther)
        return
    }
    if id >= len(DB.Data) {
        writer.WriteHeader(http.StatusNotFound)
        writer.Write([]byte("gibbet nischt"))
        return
    }

    serveTemplate("post.tmpl", writer, map[int]interface{} {id: DB.Data[id]})
}

func PostCreator(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    writer.WriteHeader(http.StatusOK)
    serveTemplate("post_creator.tmpl", writer, DB.Data)
}

func PostDeleter(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)
    if len(vars) == 0 {
        freshlog.Warn.Print("somehow passed empty vars to deleter")
        writer.WriteHeader(http.StatusInternalServerError)
    }

    id,err := strconv.Atoi(vars["id"])
    if err != nil {
        freshlog.Warn.Print("somehow passed invalid  vars to deleter, ", err)
        writer.WriteHeader(http.StatusInternalServerError)
    }
    writer.WriteHeader(http.StatusOK)
    serveTemplate("post_deleter.tmpl", writer, DB.Data[id])
}

func serveTemplate(tmpl string, writer http.ResponseWriter, data interface{}) {
    templ_file,err := ioutil.ReadFile(CONFIG.Template_path + "/" + tmpl)
    if err != nil {
        freshlog.Error.Print("failed to read template file: ", err)
        writer.WriteHeader(http.StatusInternalServerError)
    }

    templ,err := template.New(tmpl).Parse(string(templ_file))
    if err != nil{
        freshlog.Error.Print("failed to parse template: ", err)
        writer.WriteHeader(http.StatusInternalServerError)
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


    //multiplex
    r := mux.NewRouter()
    r.HandleFunc("/", MainHandler)
    r.HandleFunc("/posts", PostsViewer)
    r.HandleFunc("/posts/{id:[0-9]*}", PostViewer)
    r.HandleFunc("/posts/{id:[0-9]*}/delete", PostDeleter)
    r.HandleFunc("/posts/create", PostCreator)

    //do shit
    freshlog.Fatal.Fatal(http.ListenAndServe(CONFIG.IP + ":" + strconv.Itoa(CONFIG.Port), r))
}
