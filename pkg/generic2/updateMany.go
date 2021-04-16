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

dataProvider.updateMany('posts', {
    ids: [123, 234],
    data: { views: 0 },
})
.then(response => console.log(response));
// {
//     data: [123, 234]
// }

    updateMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        return httpClient(`${apiUrl}/${resource}?${stringify(query)}`, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },
*/

package generic
