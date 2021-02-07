/*
getList 	GET http://my.api.url/posts?sort=["title","ASC"]&range=[0, 24]&filter={"title":"bar"}
getOne 	GET http://my.api.url/posts/123
getMany 	GET http://my.api.url/posts?filter={"id":[123,456,789]}
getManyReference 	GET http://my.api.url/posts?filter={"author_id":345}
create 	POST http://my.api.url/posts
update 	PUT http://my.api.url/posts/123
updateMany 	Multiple calls to PUT http://my.api.url/posts/123
delete 	DELETE http://my.api.url/posts/123
deleteMany 	Multiple calls to DELETE http://my.api.url/posts/123

dataProvider.getOne('posts', { id: 123 })
.then(response => console.log(response));
// {
//     data: { id: 123, title: "hello, world" }
// }

    getOne: (resource, params) =>
        httpClient(`${apiUrl}/${resource}/${params.id}`).then(({ json }) => ({
            data: json,
        })),
*/

package generic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOneSample(id string) (error, map[string]interface{}) {
	m := map[string]interface{}{
		"id":    id,
		"title": "hello, world",
	}

	return nil, m
}

func GetOne(f func(string) (error, map[string]interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")

		err, data := f(id)
		if err != nil {
			Fail(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
