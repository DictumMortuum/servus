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

dataProvider.create('posts', { data: { title: "hello, world" } })
.then(response => console.log(response));
// {
//     data: { id: 450, title: "hello, world" }
// }

    create: (resource, params) =>
        httpClient(`${apiUrl}/${resource}`, {
            method: 'POST',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({
            data: { ...params.data, id: json.id },
        })),
*/

package generic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateSample(map[string]interface{}) (error, map[string]interface{}) {
	m := map[string]interface{}{
		"id":    450,
		"title": "hello, world",
	}

	return nil, m
}

func Create(f func(map[string]interface{}) (error, map[string]interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		var body map[string]interface{}
		c.Bind(&body)

		err, resp := f(body)
		if err != nil {
			Fail(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": resp,
		})
	}
}
