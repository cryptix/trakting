package main

import (
	"fmt"

	"github.com/cryptix/gocrayons"
)

func main() {
	// Create Basic Auth
	auth := gocrayons.BasicAuth{"username", "password"}
	// Create New Api with our auth
	api := gocrayons.Api("http://your-api-url.com/api/", &auth)

	// Maybe some payload to send along with the request?
	payload := map[string]interface{}{"Key": "Value1"}
	// Perform a GET request
	// URL Requested: http://your-api-url.com/api/users
	resp, err := api.Res("users").Get(nil)
	check(err)
	fmt.Printf("%+v\n", resp.Response)

	// Get Single Item
	resp, err = api.Res("users").Id(1).Get(nil)
	check(err)
	fmt.Printf("%+v\n", resp.Response)

	// Perform a GET request with Querystring
	querystring := map[string]string{"page": "100", "per_page": "1000"}
	// URL Requested: http://your-api-url.com/api/users/123/items?page=100&per_page=1000
	resp, err = api.Res("users").Id(123).Res("items").Get(querystring)
	check(err)
	fmt.Printf("%+v\n", resp.Response)

	// Or perform a POST Request
	// URL Requested: http://your-api-url.com/api/items/123 with payload as json Data
	resp, err = api.Res("items").Id(123).Post(payload)
	check(err)
	fmt.Printf("%+v\n", resp.Response)

	// Users endpoint
	users := api.Res("users")

	for i := 0; i < 10; i++ {
		// Create a new pointer to response Struct
		// user := new(respStruct)
		// Get user with id i into the newly created response struct
		resp, err := users.Id(i).Get(nil)
		check(err)
		fmt.Printf("%+v\n", resp.Response)
	}

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
