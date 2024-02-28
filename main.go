package main

import (
	"net/http"
	"strconv"

	_ "hasura/demo-blog-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/authors", getAuthors)
	router.POST("/author", addAuthor)
	router.DELETE("/author/:authorId", deleteAuthor)

	router.GET("/blogs/", getBlogsByAuthor) // author should be a query param
	router.POST("/blog", postAddBlog)
	router.PATCH("/blog/:blogId", patchUpdateBlog)
	router.PUT("/blog/", putUpdateBlog)
	router.PATCH("/blog/like", patchLikeBlog)
	router.PATCH("/blog/dislike", patchDislikeBlog)
	router.DELETE("/blog/:blogId", deleteBlog)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("localhost:9090")
}

// @Summary get all authors
// @Schemes
// @Description get all authors
// @Tags author
// @Accept json
// @Produce json
// @Success 200 {array} author
// @Router /authors [get]
func getAuthors(c *gin.Context) {
	undeletedAuthors := []author{}

	for _, element := range authorsDb {
		if !element.Deleted {
			undeletedAuthors = append(undeletedAuthors, element)
		}
	}

	c.JSON(http.StatusOK, undeletedAuthors)
}

// @Summary create author
// @Schemes
// @Description create a new author
// @Tags author
// @Param author body author true "new author to be added"
// @Produce json
// @Success 200 {object} author
// @Router /author [post]
func addAuthor(c *gin.Context) {
	var newAuthor author

	if err := c.BindJSON(&newAuthor); err != nil {
		return
	}

	newAuthor.ID = len(authorsDb)
	authorsDb = append(authorsDb, newAuthor)
	c.JSON(http.StatusCreated, newAuthor)
}

// @Summary delete author
// @Schemes
// @Description delete an author
// @Tags author
// @Produce json
// @Param authorId  path  string  true  "delete author"
// @Success 200 {object} author
// @Router /author/{authorId} [delete]
func deleteAuthor(c *gin.Context) {
	authorId, err := strconv.Atoi(c.Param("authorId"))
	if err != nil {
		return
	}
	if authorId > (len(authorsDb)+1) || authorId < 1 {
		return
	}

	authorsDb[authorId-1].Deleted = true

	c.JSON(http.StatusOK, authorsDb[authorId-1])
}

// @Summary get blogs by author
// @Schemes
// @Description get blogs by author
// @Tags blog
// @Produce json
// @Param authorId  query  string  true  "id of the author"
// @Success 200 {array} blog
// @Router /blogs/ [get]
func getBlogsByAuthor(c *gin.Context) {
	authorIdStr, ok := c.GetQuery("authorId")
	if !ok {
		return
	}
	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		return
	}
	if authorId > (len(authorsDb)+1) || authorId < 1 {
		return
	}

	blogs := []blog{}

	for _, element := range blogsDb {
		if element.Author.ID == authorId && !element.Deleted {
			blogs = append(blogs, element)
		}
	}

	c.JSON(http.StatusOK, blogs)
}

// @Summary create a new blog post
// @Schemes
// @Description create a new blog post
// @Tags blog
// @Produce json
// @Param blog  body  blogRequestDto  true  "post to add"
// @Success 200 {object} blog
// @Router /blog/ [post]
func postAddBlog(c *gin.Context) {
	var newBlogDto blogRequestDto

	if err := c.BindJSON(&newBlogDto); err != nil {
		return
	}

	authorId := newBlogDto.AuthorId
	if authorId > (len(authorsDb)+1) || authorId < 1 {
		return
	}

	newBlogDto.ID = len(blogsDb)

	newBlog := blog{
		newBlogDto.ID,
		newBlogDto.Title,
		newBlogDto.Text,
		newBlogDto.Summary,
		authorsDb[authorId-1],
		0,
		0,
		false,
	}

	blogsDb = append(blogsDb, newBlog)
	c.JSON(http.StatusCreated, newBlog)
}

// @Summary update a  blog post
// @Schemes
// @Description update a blog post
// @Tags blog
// @Produce json
// @Param blog  body  blogRequestDto  true  "blog to update"
// @Success 200 {object} blog
// @Router /blog/ [put]
func putUpdateBlog(c *gin.Context) {
	var newBlogDto blogRequestDto

	if err := c.BindJSON(&newBlogDto); err != nil {
		return
	}

	if newBlogDto.ID > (len(blogsDb)+1) || newBlogDto.ID < 1 {
		return
	}

	blogsDb[newBlogDto.ID-1].Title = newBlogDto.Title
	blogsDb[newBlogDto.ID-1].Text = newBlogDto.Text
	blogsDb[newBlogDto.ID-1].Summary = newBlogDto.Summary
	blogsDb[newBlogDto.ID-1].Author = authorsDb[newBlogDto.AuthorId-1]
	blogsDb[newBlogDto.ID-1].Likes = newBlogDto.Likes
	blogsDb[newBlogDto.ID-1].Dislikes = newBlogDto.Dislikes

	c.JSON(http.StatusOK, blogsDb[newBlogDto.ID-1])
}

// @Summary update a blog post
// @Schemes
// @Description update a blog post
// @Tags blog
// @Produce json
// @Param blogId  path  string  true  "blog to update"
// @Param title  query  string  false  "title to update"
// @Param text  query  string  false  "text to update"
// @Param summary  query  string  false  "summary to update"
// @Success 200 {object} blog
// @Router /blog/{blogId} [patch]
func patchUpdateBlog(c *gin.Context) {
	blogId, err := strconv.Atoi(c.Param("blogId"))
	if err != nil {
		return
	}
	if blogId > (len(blogsDb)+1) || blogId < 1 {
		return
	}
	title, ok := c.GetQuery("title")
	if ok {
		blogsDb[blogId-1].Title = title
	}
	text, ok := c.GetQuery("text")
	if ok {
		blogsDb[blogId-1].Text = text
	}
	summary, ok := c.GetQuery("summary")
	if ok {
		blogsDb[blogId-1].Summary = summary
	}
	c.JSON(http.StatusOK, blogsDb[blogId-1])
}

// @Summary like a blog post
// @Schemes
// @Description like a blog post
// @Tags blog
// @Produce json
// @Param blogId  query  string  true  "blog to like"
// @Success 200 {object} blog
// @Router /blog/like [patch]
func patchLikeBlog(c *gin.Context) {
	blogIdStr, ok := c.GetQuery("blogId")
	if ok {
		return
	}
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		return
	}
	if blogId > (len(blogsDb)+1) || blogId < 1 {
		return
	}
	blogsDb[blogId-1].Likes += 1
	c.JSON(http.StatusOK, blogsDb[blogId-1])
}

// @Summary dislike a blog post
// @Schemes
// @Description dislike a blog post
// @Tags blog
// @Produce json
// @Param blogId  query  string  true  "blog to dislike"
// @Success 200 {object} blog
// @Router /blog/dislike [patch]
func patchDislikeBlog(c *gin.Context) {
	blogIdStr, ok := c.GetQuery("blogId")
	if ok {
		return
	}
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		return
	}
	if blogId > (len(blogsDb)+1) || blogId < 1 {
		return
	}
	blogsDb[blogId-1].Likes -= 1
	c.JSON(http.StatusOK, blogsDb[blogId-1])
}

// @Summary delete a blog post
// @Schemes
// @Description delete a blog post
// @Tags blog
// @Produce json
// @Param blogId  path  string  true  "blog to delete"
// @Success 200 {object} blog
// @Router /blog/{blogId} [delete]
func deleteBlog(c *gin.Context) {
	blogId, err := strconv.Atoi(c.Param("blogId"))
	if err != nil {
		return
	}
	if blogId > (len(blogsDb)+1) || blogId < 1 {
		return
	}

	blogsDb[blogId-1].Deleted = true

	c.JSON(http.StatusOK, blogsDb[blogId-1])
}

type author struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Deleted bool   `json:"-"`
}

var authorsDb = []author{
	{1, "Critic", "critic@demoblog.api", false},
	{2, "Philospoher", "philospoher@demoblog.api", false},
	{3, "Artist", "artist@demoblog.api", false},
	{4, "Dreamer", "dreamer@demoblog.api", false},
	{5, "Poet", "poet@demoblog.api", false},
}

type blog struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Summary  string `json:"summary"`
	Author   author `json:"author"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Deleted  bool   `json:"-"`
}

type blogRequestDto struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Summary  string `json:"summary"`
	AuthorId int    `json:"authorId"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

var blogsDb = []blog{
	{
		1,
		"Reflections minus traumas",
		"Someday, we'll complain. Not that life was difficult and we were victims, but that life was beautiful, and we didn't pay attention",
		"",
		authorsDb[1],
		10,
		2,
		false,
	},
	{
		2,
		"Alarmingly mid",
		"I see a sea of people, drowning in mediocrity because they've accepted 'good enough' to be their standard. Oh my fellow human how I yearn to let you know! You were meant to be nothing short of excellent",
		"",
		authorsDb[0],
		7,
		0,
		false,
	},
	{
		3,
		"Far away",
		`Under an open sky I lay, 
thinking of you - dare I say, 
billions of heavenly jewels in space, 
and yet - not a single one more beautiful than your face, 
with distances so vast they're measured in light years, 
they pale in comparison to our separation I fear`,
		"",
		authorsDb[4],
		13,
		1,
		false,
	},
	{
		4,
		"Colors",
		"Can you imagine? A monotone world? A world without color. A world where you've seen everything? There is no new, just the same. Everyday is every other day, with different shades of grey? Bring some color into your life my friend, bring some color into your life",
		"",
		authorsDb[2],
		10,
		2,
		false,
	},
}
