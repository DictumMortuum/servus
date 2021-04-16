/*
getList 	GET http://my.api.url/posts?sort=["title","ASC"]&range=[0, 24]&filter={"title":"bar"}

Note: The simple REST client expects the API to include a Content-Range header in the response to getList calls.
The value must be the total number of resources in the collection. This allows react-admin to know how many pages of resources there are in total,
and build the pagination controls.
Content-Range: posts 0-24/319
If your API is on another domain as the JS code, you’ll need to whitelist this header with an Access-Control-Expose-Headers CORS header.
Access-Control-Expose-Headers: Content-Range

https://marmelab.com/react-admin/DataProviders.html#request-format
https://github.com/gin-gonic/gin
https://github.com/gin-gonic/gin/issues/1516

dataProvider.getList('posts', {
    pagination: { page: 1, perPage: 5 },
    sort: { field: 'title', order: 'ASC' },
    filter: { author_id: 12 },
})
.then(response => console.log(response));
// {
//     data: [
//         { id: 126, title: "allo?", author_id: 12 },
//         { id: 127, title: "bien le bonjour", author_id: 12 },
//         { id: 124, title: "good day sunshine", author_id: 12 },
//         { id: 123, title: "hello, world", author_id: 12 },
//         { id: 125, title: "howdy partner", author_id: 12 },
//     ],
//     total: 27
// }

    getList: (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
            sort: JSON.stringify([field, order]),
            range: JSON.stringify([(page - 1) * perPage, page * perPage - 1]),
            filter: JSON.stringify(params.filter),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        return httpClient(url).then(({ headers, json }) => ({
            data: json,
            total: parseInt(headers.get('content-range').split('/').pop(), 10),
        }));
    },
*/

package generic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetListSample(args *Args) (error, int, []map[string]interface{}) {
	m := []map[string]interface{}{
		{
			"id":    1,
			"title": "hello, world",
		}, {
			"id":    2,
			"title": "hello, world",
		}, {
			"id":    3,
			"title": "hello, world",
		}, {
			"id":    4,
			"title": "hello, world",
		}, {
			"id":    5,
			"title": "hello, world",
		}, {
			"id":    6,
			"title": "hello, world",
		}, {
			"id":    7,
			"title": "hello, world",
		}, {
			"id":    8,
			"title": "hello, world",
		}, {
			"id":    9,
			"title": "hello, world",
		},
	}

	var low, high int
	max := len(m)

	if max >= args.Range[1] {
		high = args.Range[1]
	} else {
		high = max
	}

	if high >= args.Range[0] {
		low = args.Range[0]
	} else {
		low = high
	}

	return nil, len(m), m[low:high:max]
}

func GetList(f func(*Args) (error, int, *[]interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		err, args := ParseArgs(c)
		if err != nil {
			Fail(c, err)
			return
		}

		log.Println(args)

		err, total, data := f(args)
		if err != nil {
			Fail(c, err)
			return
		}

		c.Header("Content-Range", fmt.Sprintf("%d-%d/%d", args.Range[0], args.Range[1], total))
		c.JSON(http.StatusOK, data)
	}
}
