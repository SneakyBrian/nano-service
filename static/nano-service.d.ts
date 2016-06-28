/**
 * simple synchronous http get with $get
 * Example: var result = $get(url);
 */
declare function $get(url: string): string;

/**
 * simple synchronous http post with $post
 * Example: var result = $post(url, "application/json", JSON.stringify(obj))
 */
declare function $post(url: string, bodyType: string, body: string): string;

/**
 * query parameters are available as the $query object
 * Example: var foo = $query.foo;
 * Example: var foo = $query["foo"];
 */
declare var $query: {
    [key: string]: string;
}