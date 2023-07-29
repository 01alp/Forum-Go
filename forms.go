package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gofrs/uuid"
)

func likePost(w http.ResponseWriter, r *http.Request) {
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if post_id == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes

	reactionsPosts := "postsReactions"
	postsTable := "posts"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			// fmt.Println("adding notification, recepient", fetchPostByID(database, post_id).UserId)
			fetch := fetchReactionByUserAndId(database, reactionsPosts, user.Id, post_id)
			action := "liked"
			if fetch.Value != 1 { // if not like
				if fetch.Value == -1 {
					// delete dislike in frontend and backend
					deleteRow(database, reactionsPosts, fetch.Id)
					updateTableDislikes(database, postsTable, allDislikes-1, post_id)
					action = "changed dislike to like on"
				}
				// add like in frontend and backend
				updateTableLikes(database, postsTable, allLikes+1, post_id)
				addPostsReactions(database, 1, user.Id, post_id)
			} else {
				// delete like in frontend and backend
				updateTableLikes(database, postsTable, allLikes-1, post_id)
				deleteRow(database, reactionsPosts, fetch.Id)
				action = "removed like from"
			}
			addNotification(database, "post", fetchPostByID(database, post_id).Title, post_id, action, user.Username, fetchPostByID(database, post_id).UserId)

		}
		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func dislikePost(w http.ResponseWriter, r *http.Request) {
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if post_id == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes

	reactionsPosts := "postsReactions"
	postsTable := "posts"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndId(database, reactionsPosts, user.Id, post_id)
			action := "disliked"
			if fetch.Value != -1 { // if not dislike
				if fetch.Value == 1 {
					// delete like in frontend and backend
					deleteRow(database, reactionsPosts, fetch.Id)
					updateTableLikes(database, postsTable, allLikes-1, post_id)
					action = "changed like to dislike on"

				}
				// add dislike in frontend and backend
				updateTableDislikes(database, postsTable, allDislikes+1, post_id)
				addPostsReactions(database, -1, user.Id, post_id)
			} else {
				// delete dislike in frontend and backend
				updateTableDislikes(database, postsTable, allDislikes-1, post_id)
				deleteRow(database, reactionsPosts, fetch.Id)
				action = "removed dislike from"
			}
			addNotification(database, "post", fetchPostByID(database, post_id).Title, post_id, action, user.Username, fetchPostByID(database, post_id).UserId)
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func likeComment(w http.ResponseWriter, r *http.Request) {
	comment_id, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if comment_id == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user := welcome(w, r).User
	allLikes := fetchCommentByID(database, comment_id).Likes
	allDislikes := fetchCommentByID(database, comment_id).Dislikes

	reactionsComments := "commentsReactions"
	commentsTable := "comments"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		// if comment_id > 0 && comment_id <= len(comments) {
		fetch := fetchReactionByUserAndId(database, reactionsComments, user.Id, comment_id)
		action := "liked"
		if fetch.Value != 1 { // if not like
			if fetch.Value == -1 {
				// delete dislike in frontend and backend
				deleteRow(database, reactionsComments, fetch.Id)
				updateTableDislikes(database, commentsTable, allDislikes-1, comment_id)
				action = "changed dislike to like on"

			}
			// add like in frontend and backend
			updateTableLikes(database, commentsTable, allLikes+1, comment_id)
			addCommentsReactions(database, 1, user.Id, comment_id)
		} else {
			// delete like in frontend and backend
			updateTableLikes(database, commentsTable, allLikes-1, comment_id)
			deleteRow(database, reactionsComments, fetch.Id)
			action = "removed like from"
		}
		post := fetchPostByID(database, fetchCommentByID(database, comment_id).PostId)
		addNotification(database, "comment", post.Title, post.Id, action, user.Username, fetchCommentByID(database, comment_id).UserId)

		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
		// }
	}
}

func dislikeComment(w http.ResponseWriter, r *http.Request) {
	comment_id, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if comment_id == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user := welcome(w, r).User
	allLikes := fetchCommentByID(database, comment_id).Likes
	allDislikes := fetchCommentByID(database, comment_id).Dislikes

	reactionsComments := "commentsReactions"
	commentsTable := "comments"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		// if comment_id > 0 && comment_id <= len(comments) {
		fetch := fetchReactionByUserAndId(database, reactionsComments, user.Id, comment_id)
		action := "disliked"
		if fetch.Value != -1 { // if not dislike
			if fetch.Value == 1 {
				// delete like in frontend and backend
				deleteRow(database, reactionsComments, fetch.Id)
				updateTableLikes(database, commentsTable, allLikes-1, comment_id)
				action = "changed like to dislike on"

			}
			// add dislike in frontend and backend
			updateTableDislikes(database, commentsTable, allDislikes+1, comment_id)
			addCommentsReactions(database, -1, user.Id, comment_id)
		} else {
			// delete dislike in frontend and backend
			updateTableDislikes(database, commentsTable, allDislikes-1, comment_id)
			deleteRow(database, reactionsComments, fetch.Id)
			action = "removed dislike from"
		}
		post := fetchPostByID(database, fetchCommentByID(database, comment_id).PostId)
		addNotification(database, "comment", post.Title, post.Id, action, user.Username, fetchCommentByID(database, comment_id).UserId)

		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
		// }
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	r.ParseForm()
	// Fetch image from form if exists
	file, header, err := r.FormFile("image")

	data.Message = &Message{
		Threads:     r.Form["threads"],
		ImageHeader: header,
	}

	fileName := ""
	if err == nil {

		fmt.Printf("Uploaded File: %+v ", header.Filename)
		fmt.Printf("File Size: %+v ", header.Size)
		fmt.Println("Checking file before upload...")

		if !data.Message.ValidateImage() {
			data.Post.Title = r.FormValue("title")
			data.Post.Content = r.FormValue("content")
			tmpl, err := template.ParseFiles("static/template/newPost.html", "static/template/base.html")
			if err != nil {
				createError(w, r, http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				createError(w, r, http.StatusInternalServerError)
				return
			}
			return

		}

		// upload the image file to static/template/assets/img/
		// Add uuid to the file name to avoid overwriting
		u, _ := uuid.NewV4()
		fileName = u.String() + header.Filename
		f, err := os.OpenFile("./static/template/assets/img/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	if !data.Message.ValidateThreads() {
		data.Post.Title = r.FormValue("title")
		data.Post.Content = r.FormValue("content")
		tmpl, err := template.ParseFiles("static/template/newPost.html", "static/template/base.html")
		if err != nil {
			createError(w, r, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			createError(w, r, http.StatusInternalServerError)
			return
		}
		return
	}

	addPost(database, r.FormValue("title"), fileName, r.FormValue("content"), r.Form["threads"], data.User.Id, 0, 0)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func createComment(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)

	// msg := &Message{}

	// if !msg.ValidateComment() {
	// 	data := Data{Message: msg, Post: Post{Title: r.FormValue("title"), Content: r.FormValue("content")}, Threads: fetchAllThreads(database)}
	// 	fmt.Println(data.Post)
	// 	tmpl, _ := template.ParseFiles("static/template/newPost.html", "static/template/base.html")
	// 	tmpl.Execute(w, data)
	// 	return
	// }
	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
		return
	}

	post, _ := strconv.Atoi(r.FormValue("id"))
	addComment(database, r.FormValue("comment"), post, data.User.Id, 0, 0)
	addNotification(database, "post", fetchPostByID(database, post).Title, post, "commented", data.User.Username, fetchPostByID(database, post).UserId)

	// c, err := r.Cookie("last_page")
	// if err != nil {
	// 	fmt.Println("cookie err", err)
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// }
	// fmt.Println("REDIRECT TO", c.Value)
	// http.Redirect(w, r, c.Value, http.StatusSeeOther)
	http.Redirect(w, r, "/post/id?id="+r.FormValue("id"), http.StatusSeeOther)
}

func disableNotifications(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)

	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
		return
	}

	disable(fetchActiveNotificationsByUserId(database, data.User.Id))

	c, err := r.Cookie("last_page")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	http.Redirect(w, r, c.Value, http.StatusSeeOther)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	r.ParseForm()
	// Fetch image from form if exists
	file, header, err := r.FormFile("image")

	data.Message = &Message{
		Threads:     r.Form["threads"],
		ImageHeader: header,
	}

	fileName := ""
	if err == nil {

		fmt.Printf("Uploaded File: %+v ", header.Filename)
		fmt.Printf("File Size: %+v ", header.Size)
		fmt.Println("Checking file before upload...")

		if !data.Message.ValidateImage() {
			data.Post.Title = r.FormValue("title")
			data.Post.Content = r.FormValue("content")
			tmpl, err := template.ParseFiles("static/template/editPost.html", "static/template/base.html")
			if err != nil {
				createError(w, r, http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				createError(w, r, http.StatusInternalServerError)
				return
			}
			return

		}

		// upload the image file to static/template/assets/img/
		// Add uuid to the file name to avoid overwriting
		u, _ := uuid.NewV4()
		fileName = u.String() + header.Filename
		f, err := os.OpenFile("./static/template/assets/img/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	err = updatePostByID(database, id, r.FormValue("title"), fileName, r.FormValue("content"), r.Form["threads"])
	if err != nil {
		createError(w, r, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/id?id="+r.FormValue("id"), http.StatusSeeOther)
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)

	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
		return
	}
	// data.Message = &Message{
	// 	Threads:     r.Form["threads"],
	// 	ImageHeader: header,
	// }

	id, _ := strconv.Atoi(r.FormValue("comment_id"))
	err := updateCommentByID(database, id, r.FormValue("content"))
	if err != nil {
		fmt.Println(err)
		createError(w, r, http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("last_page")
	// if err != nil {
	// 	http.Redirect(w, r, "/post/id?id="+r.FormValue(""), http.StatusSeeOther)
	// }
	http.Redirect(w, r, c.Value, http.StatusSeeOther)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	id, _ := strconv.Atoi(r.FormValue("id"))

	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
		return
	}
	if data.User.Id != fetchPostByID(database, id).UserId { // wrong user
		createError(w, r, http.StatusBadRequest)
		return
	}
	// data.Message = &Message{
	// 	Threads:     r.Form["threads"],
	// 	ImageHeader: header,
	// }

	err := deleteRow(database, "posts", id)
	if err != nil {
		fmt.Println(err)
		createError(w, r, http.StatusInternalServerError)
		return
	}

	// c, err := r.Cookie("last_page")
	// if err != nil {
	 	// http.Redirect(w, r, "/post/id?id="+r.FormValue(""), http.StatusSeeOther)
	// }
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	id, _ := strconv.Atoi(r.FormValue("comment_id"))

	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
		return
	}
	if data.User.Id != fetchCommentByID(database, id).UserId { // wrong user
		createError(w, r, http.StatusBadRequest)
		return
	}
	// data.Message = &Message{
	// 	Threads:     r.Form["threads"],
	// 	ImageHeader: header,
	// }

	err := deleteRow(database, "comments", id)
	if err != nil {
		fmt.Println(err)
		createError(w, r, http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("last_page")
	// if err != nil {
	// 	http.Redirect(w, r, "/post/id?id="+r.FormValue(""), http.StatusSeeOther)
	// }
	http.Redirect(w, r, c.Value, http.StatusSeeOther)
}
