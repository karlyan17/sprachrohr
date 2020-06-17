// sprachrohr.go
package main


import(
    "github.com/karlyan17/jimbob"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "github.com/gorilla/securecookie"
    "net/http"
    "strconv"
    "html/template"
    "golang.org/x/crypto/bcrypt"
    "sprachrohr/post"
    "sprachrohr/freshlog"
    "sprachrohr/config"
)

var CONFIG  config.Config
var POSTS   jimbob.Bucket
var USERS   jimbob.Bucket
var COOK    *sessions.CookieStore

func MainHandler(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    http.Redirect(writer, request, "/posts", http.StatusSeeOther)
}

func AuthHandler(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

    switch request.Method {
    case http.MethodPost:
        user_name := request.PostFormValue("user")
        freshlog.Debug.Print("user trying to log in: ", user_name)
        plain_pw := []byte(request.PostFormValue("pw"))

		auth := false
		if len(USERS.Data) == 0 {
            freshlog.Debug.Print("no users in db creating first")
            pw,err := bcrypt.GenerateFromPassword(plain_pw , bcrypt.DefaultCost)
            if err != nil {
                freshlog.Error.Print("unable to hash password: ", err)
            }
            user := map[string]interface{}{"user": user_name, "pw": pw}
			id, err := USERS.Post(user)
			if err != nil {
				freshlog.Warn.Print("committing to database failed: ", err)
			}
            freshlog.Debug.Print("created user: ", user_name, " id: ", id)
			auth = true
		} else {
            for _,user_data := range(USERS.Data) {
                user_obj := user_data.(map [string]interface{})
                if user_obj["user"] == user_name {
                    err := bcrypt.CompareHashAndPassword([]byte(user_obj["pw"].(string)), []byte(plain_pw))
                    if err == nil {
                        auth = true
                    }
                }
            }
        }

        if !auth {
            freshlog.Warn.Print("authentification failed: ", user_name)
        }
		session, err := COOK.Get(request, "sprachrohr-sess")
		if err != nil {
			freshlog.Warn.Print("error decoding session: ", err)
		}
		session.Values["auth"] = auth
		session.Save(request, writer)

    case http.MethodDelete:
		session, err := COOK.Get(request, "sprachrohr-sess")
		if err != nil {
			freshlog.Warn.Print("error decoding session: ", err)
		}
		session.Values["auth"] = false
		session.Save(request, writer)
	}

	http.Redirect(writer, request, "/posts", http.StatusSeeOther)
	return
}

func PostsViewer(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

	session, err := COOK.Get(request, "sprachrohr-sess")
    if err != nil {
        freshlog.Warn.Print("error decoding session: ", err)
    }

    data := POSTS.Data

    serveTemplate("posts.tmpl", writer, map[string]interface{} {"data": data, "session": session})
}

func PostViewer(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

	session, err := COOK.Get(request, "sprachrohr-sess")
    if err != nil {
        freshlog.Warn.Print("error decoding session: ", err)
    }

    if len(vars) == 0 {
        http.Redirect(writer, request, "/posts", http.StatusSeeOther)
        return
    }

    id,err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Redirect(writer, request, "/posts", http.StatusSeeOther)
        return
    }
    if POSTS.Data[id] == nil {
        //TODO log not found error
        http.Error(writer, "gibbet nischt", http.StatusNotFound)
        return
    }

    data := map[int]interface{} {id: POSTS.Data[id]}

    serveTemplate("post.tmpl", writer, map[string]interface{} {"data": data, "session": session})
}

func PostCreator(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

	session, err := COOK.Get(request, "sprachrohr-sess")
    if err != nil {
        freshlog.Warn.Print("error decoding session: ", err)
    }

    if auth, ok := session.Values["auth"].(bool); !ok || !auth {
		freshlog.Warn.Print("forbidden request: ", err)
        http.Error(writer, "fuck off", http.StatusForbidden)
        return
    }

    switch request.Method {
    case http.MethodGet:
        writer.WriteHeader(http.StatusOK)
        data := POSTS.Data
        serveTemplate("post_creator.tmpl", writer, map[string]interface{} {"data": data, "session": session})
    case http.MethodPost:
        err := request.ParseForm()
        if err != nil {
            freshlog.Warn.Print("error parsing Form: ", err)
            http.Redirect(writer, request, request.RequestURI, http.StatusNotModified)
            return
        }
        title := request.PostFormValue("title")
        freshlog.Debug.Print("Title: ", title)
        body := request.PostFormValue("body")
        freshlog.Debug.Print("Body: ", body)
        if title == "" || body == "" {
            freshlog.Warn.Print("neither Title nor Body can be empty!")
            http.Redirect(writer, request, request.RequestURI, http.StatusNotModified)
            return
        }
        new_post := post.NewPost(title, body)
        id, err := POSTS.Post(new_post)
        if err != nil {
            freshlog.Warn.Print("committing to database failed: ", err)
            http.Redirect(writer, request, request.RequestURI, http.StatusNotModified)
            return
        }
        http.Redirect(writer, request, "/posts/" + strconv.Itoa(id), http.StatusFound)
        return
    }

}

func PostDeleter(writer http.ResponseWriter, request *http.Request) {
    freshlog.Debug.Print("request is: ", request)

    vars := mux.Vars(request)
    freshlog.Debug.Print("vars are: ", vars)

	session, err := COOK.Get(request, "sprachrohr-sess")

    if auth, ok := session.Values["auth"].(bool); !ok || !auth {
		freshlog.Warn.Print("forbidden request: ", err)
        http.Error(writer, "fuck off", http.StatusForbidden)
        return
    }

    if len(vars) == 0 {
        freshlog.Warn.Print("somehow passed empty vars to deleter")
        writer.WriteHeader(http.StatusInternalServerError)
    }

    id,err := strconv.Atoi(vars["id"])
    if err != nil {
        //TODO make error
        http.Redirect(writer, request, "/posts", http.StatusSeeOther)
        return
    }
    if POSTS.Data[id] == nil {
        //TODO make error
        http.Error(writer, "gibbet nischt", http.StatusNotFound)
        return
    }

    switch request.Method {
    case http.MethodGet:
        data := map[int]interface{} {id: POSTS.Data[id]}
        writer.WriteHeader(http.StatusOK)
        serveTemplate("post_deleter.tmpl", writer, map[string]interface{} {"data": data, "session": session})
        return
    case http.MethodPost:
        if err != nil {
            freshlog.Warn.Print("error converting ID to string: ", err)
            http.Redirect(writer, request, request.RequestURI, http.StatusNotModified)
            return
        }
        err = POSTS.Delete(id)
        if err != nil {
            freshlog.Warn.Print("committing to database failed: ", err)
            http.Redirect(writer, request, request.RequestURI, http.StatusNotModified)
            return
        }
        freshlog.Debug.Print("Successfully deleted: ", id)
        http.Redirect(writer, request, "/posts", http.StatusFound)
        return
    }
}

func serveTemplate(tmpl string, writer http.ResponseWriter, data interface{}) {
    //TODO return error
    templ,err := template.ParseFiles(CONFIG.Template_path + "/__first.tmpl", CONFIG.Template_path + "/" +  tmpl, CONFIG.Template_path + "/__last.tmpl")
    if err != nil {
        freshlog.Error.Print("failed to read template file: ", err)
        writer.WriteHeader(http.StatusInternalServerError)
    }

    err = templ.ExecuteTemplate(writer, "__first.tmpl", data)
    if err != nil {
        freshlog.Error.Print("template error ", err)
        return
    }
    err = templ.ExecuteTemplate(writer, tmpl,  data)
    if err != nil {
        freshlog.Error.Print("template error ", err)
        return
    }

    err = templ.ExecuteTemplate(writer, "__last.tmpl", data)
    if err != nil {
        freshlog.Error.Print("template error ", err)
        return
    }
}

func main() {
    //
    CONFIG = config.ParseFlags()
    freshlog.SetLogLevel(CONFIG.Log_Level)
    freshlog.Debug.Print(CONFIG)

    freshlog.Debug.Print("opening jimbob bucket posts")
    var err error
    POSTS,err = jimbob.OpenBucket(CONFIG.DB_path + "/posts")
    if err != nil {
        freshlog.Fatal.Fatal("could not open jimbob Bucket: ",err)
    }

    freshlog.Debug.Print("opening jimbob bucket users")
    USERS,err = jimbob.OpenBucket(CONFIG.DB_path + "/users")
    if err != nil {
        freshlog.Fatal.Fatal("could not open jimbob Bucket: ",err)
    }

    freshlog.Debug.Print("opening cookie jar")
    COOK = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

    //multiplex
    r := mux.NewRouter()
    static_dir := http.Dir("static")
    static_handler := http.FileServer(static_dir)
    r.HandleFunc("/", MainHandler).Methods("GET")
    r.PathPrefix("/static/{.+}").Handler(http.StripPrefix("/static/", static_handler))
    r.HandleFunc("/auth", AuthHandler).Methods("POST","DELETE")
    r.HandleFunc("/posts", PostsViewer).Methods("GET")
    r.HandleFunc("/posts/{id:[0-9]*}", PostViewer).Methods("GET")
    r.HandleFunc("/posts/{id:[0-9]*}/delete", PostDeleter).Methods("GET","POST")
    r.HandleFunc("/posts/create", PostCreator).Methods("GET","POST")

    //do shit
    freshlog.Fatal.Fatal(http.ListenAndServe(CONFIG.IP + ":" + strconv.Itoa(CONFIG.Port), r))
}
