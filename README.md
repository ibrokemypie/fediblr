# README

Posts the latest image from a tumblr blog to fedi.

Run regularly with cron if you wish, it has a check to avoid posting the same thing repeatedly if nothing new has been posted on the blog.

The configuration (config.toml) is generated interactively, but the values are as follows

```
visiblity (the post visiblity on fedi, defaults to unlisted)

fediInstance (the url of the target instance. includes protocol)

fediToken (the authorization key for oauth with fedi)

lastImage (the file name of the last image posted. used to prevent duplicate posts.)

tumblrKey (the oauth consumer key for tumblr)

tumblrUser (the tumblr blog username or url to retrieve posts from)
```
