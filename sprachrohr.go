// sprachrohr.go
package main


import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
)

// variables
var environ []string
var response_body string
var blog_path string = "/home/nurgling/blog"

type Post struct {
        Date string
        Title string
        Body string
        Comments [5]string
}

var page [10]Post

func main() {
	env_var := make(map[string]string)
	environ = os.Environ()
	for _,n := range(environ) {
		split := strings.SplitN(n, "=", 2)
		env_var[split[0]] = split[1]
		//response_body = response_body + n + "\n"
	}
	//fmt.Println(env_var)
	postlist,_ := ioutil.ReadDir("/home/nurgling/blog")

	for i,p := range(page) {
		//fmt.Println(postlist[i].Name())
		content,_ := ioutil.ReadFile(blog_path + "/" + postlist[i].Name())
		json.Unmarshal(content, &p)
		post_body := fmt.Sprintf("<h3>%v</h3>%v\n<p>%v</p>\n<p>Comments:<ul>\n", p.Title, p.Date,  p.Body)
		for _,c := range(p.Comments) {
			if  c != "" {
				post_body += fmt.Sprintf("<li>%v\n", c)
			}
		}
		post_body += "</ul><br>\n"
		response_body = post_body + response_body
		if i == len(postlist) - 1 {
			break
		}
	}
	response_body = "<!doctype html>\n<html>\n<header><title>SprachRohr Blog</title></header>\n<body>" + response_body + "</body>\n</html>\n"

	fmt.Printf("HTTP/1.1 200 OK\r\nServer: nurgling/0.1\r\n")
	fmt.Printf("Content-Type: text/html\r\nContent-Length: %v\r\n\r\n" + response_body, len(response_body))
}
