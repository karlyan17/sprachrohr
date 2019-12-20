// post.go
package post

import(
    "time"
)

type Post struct {
    Created time.Time
    Title string
    Body string
    Comments [5]string
}


func NewPost(title string, body string) Post {
    post := Post{
        Created: time.Now(),
        Title: title,
        Body: body,
        Comments: [5]string{},
    }

    return post
}


func (post Post) CommentOn(comment string) Post {
    var new_comments [5]string
    var new_post Post

    new_comments = shiftRComment(post.Comments)
    new_comments[0] = comment
    new_post = post
    new_post.Comments = new_comments
    return new_post
}


func shiftRComment(array [5]string) [5]string {
    var new_array [5]string

    for i,comment := range(array) {
        if i < len(array) {
            new_array[i+1] = comment
        } else {
            new_array[0] = comment
        }
    }

    return new_array
}
