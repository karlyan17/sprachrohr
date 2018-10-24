// sprachrohr.go
package main


import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
	"strconv"
	"time"
)

// variables
var environ []string
var args string
var response_body string
var blog_path string = "/home/nurgling/blog"
var env_var map[string]string
var query_var map[string]string

type Post struct {
        Date string
        Title string
        Body string
        Comments [5]string
}

var page [10]Post

func buildPage(pagenum int) string {
	var response_body string
	postlist,_ := ioutil.ReadDir("/home/nurgling/blog")
	for i,p := range(page) {
		//fmt.Println(postlist[i].Name())
		content,_ := ioutil.ReadFile(blog_path + "/" + postlist[i + (pagenum - 1) * len(page)].Name())
		json.Unmarshal(content, &p)
		epoch,_ := strconv.Atoi(p.Date)
		post_time := time.Unix(int64(epoch), 0)
		post_time_string := fmt.Sprintf("%d-%02d-%02d", post_time.Year(), post_time.Month(), post_time.Day())
		post_body := fmt.Sprintf("<h3>%v<a href=\"?perma=%v\"> &gt;&gt;</a></h3>", p.Title, p.Date)
		post_body += fmt.Sprintf("<p>%v<br>\n%v<br>\nComments:</p>\n<ul>\n", post_time_string,  p.Body)
		for _,c := range(p.Comments) {
			if  c != "" {
				post_body += fmt.Sprintln("<li>" + c + "\n")
			}
		}
		post_body += "</ul><br>\n"
		post_body += "<form action=\"" + env_var["REQUEST_URI"] + "\" method=\"POST\">"
		post_body += "<input type=\"hidden\" name=\"id\" value=\"" + p.Date + "\" />"
		post_body += "<input type=\"submit\" value=\"Comment\" style=\"height:40px; width:80px\">"
		post_body += "<textarea name=\"c\" rows=\"2\" cols=\"50\" size=\"1000\">"
		post_body += "</textarea>"
		post_body += "</form>"
		response_body = post_body + response_body
		if i == len(postlist) - 1 {
			break
		}
	}
	return response_body
}
func buildPost(id string) string {
	var response_body string
	var p Post
	content,err := ioutil.ReadFile(blog_path + "/" + id)
	if err == nil {
		json.Unmarshal(content, &p)
		epoch,_ := strconv.Atoi(p.Date)
		post_time := time.Unix(int64(epoch), 0)
		post_time_string := fmt.Sprintf("%d-%02d-%02d", post_time.Year(), post_time.Month(), post_time.Day())
		post_body := fmt.Sprintf("<h3>%v<a href=\"?perma=%v\"> &gt;&gt;</a></h3>", p.Title, p.Date)
		post_body += fmt.Sprintf("<p>%v<br>\n%v<br>\nComments:</p>\n<ul>\n", post_time_string,  p.Body)
		for _,c := range(p.Comments) {
			if  c != "" {
				post_body += fmt.Sprintln("<li>" + c + "\n")
			}
		}
		post_body += "</ul><br>\n"
		post_body += "<form action=\"" + env_var["REQUEST_URI"] + "\" method=\"POST\">"
		post_body += "<input type=\"hidden\" name=\"id\" value=\"" + p.Date + "\" />"
		post_body += "<input type=\"submit\" value=\"Comment\" style=\"height:40px; width:80px\">"
		post_body += "<textarea name=\"c\" rows=\"2\" cols=\"50\" size=\"1000\">"
		post_body += "</textarea>"
		post_body += "</form>"
		response_body = post_body + response_body
	} else {
		response_body = "post " + id + " not found"
	}
	return response_body
}

func updateComment(id, comment string) {
	var post Post
	content,err := ioutil.ReadFile(blog_path + "/" + id)
	json.Unmarshal(content, &post)
	if err != nil {
		response_body += "<h2>go fuck yourself</h2>"
	} else {
		if post.Comments[0] != "" {
			for j:=len(post.Comments)-1;j>0;j-- {
				post.Comments[j] = post.Comments[j-1]
			}
		}
		post.Comments[0] = comment
		json_bytes,_ := json.Marshal(post)
		ioutil.WriteFile(blog_path + "/" + id, json_bytes, 0644)
	}
}

func main() {
	env_var = make(map[string]string)
	query_var = make(map[string]string)
	environ = os.Environ()
	for _,n := range(environ) {
		split := strings.SplitN(n, "=", 2)
		env_var[split[0]] = split[1]
	}
	if env_var["QUERY_STRING"] != "" {
		split_query := strings.Split(env_var["QUERY_STRING"], "&")
		for _,n := range(split_query) {
			split := strings.SplitN(n, "=", 2)
			query_var[split[0]] = split[1]
		}
	}
	response_body = "<!doctype html>\n<html><meta charset=\"utf-8\">\n"
	response_body += "<header><title>SprachRohr Blog</title></header>\n<body>\n" 
	response_body += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"
	response_body += "<img src=\"/podge.png\" width=\"134\" height=\"90\">\n"
	response_body += "<h1>SprachRohr Blog</h1>\n"
	if env_var["REQUEST_METHOD"] == "POST" {
		if len(os.Args) == 2 && os.Args[1] != "" {
			args,_ = url.QueryUnescape(os.Args[1])
			arg_var := make(map[string]string)
			split_args := strings.Split(args, "&")
			for _,n := range(split_args) {
				split := strings.SplitN(n, "=", 2)
				arg_var[split[0]] = split[1]
			}
			if arg_var["c"] != "" && arg_var["id"] != "" {
				clean_c := strings.Replace(arg_var["c"], "&", "&amp; ", -1)
				clean_c = strings.Replace(clean_c, ">", "&gt; ", -1)
				clean_c = strings.Replace(clean_c, "<", "&lt; ", -1)
				clean_c = strings.Replace(clean_c, "\"", "&quot; ", -1)
				clean_c = strings.Replace(clean_c, "'", "&apos; ", -1)
				if len(clean_c) > 1000 {
					clean_c = string([]rune(clean_c)[0:1000])
				}
				updateComment(arg_var["id"], clean_c)
			}
		}
	}
	if id := query_var["perma"]; id != "" {
		response_body += buildPost(id)
	} else {
		response_body += buildPage(1)
	}
	response_body += "<p>For bugs, ideas, suggestion and other spam: karlyan.kamerer (at) gmail.com </p>"
	response_body += "</body>\n</html>\n"

	fmt.Printf("HTTP/1.1 200 OK\r\nServer: nurgling/0.1\r\n")
	fmt.Printf("Content-Type: text/html; charset=utf-8\r\nContent-Length: %v\r\n\r\n" + response_body, len(response_body))
}
