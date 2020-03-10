package tumblr

import "fmt"

func GetTumblrKey() string {
	fmt.Println("Register a tumblr application at https://www.tumblr.com/oauth/apps")
	fmt.Println("The paste the OAuth consumer key from that page and press enter.")

	var tumblrKey string
	fmt.Scanln(&tumblrKey)

	return tumblrKey
}

func GetTumblrUser() string {
	fmt.Println("Paste the tumblr username or blog url (tumblr.com/user) and press enter")

	var tumblrUser string
	fmt.Scanln(&tumblrUser)

	return tumblrUser
}
