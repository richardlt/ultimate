package main

// User struct
//ultimate:sql
type User struct {
	id        int64  `ultimate:"id,criteria"`
	name      string `ultimate:"name,criteria"`
	ignoreOne string
	IgnoreTwo string `json:"ignore-two"`
}

//  unknown

// Group struct
//ultimate:mongo
type Group struct {
	id    int64   `ultimate:"id"`
	name  string  `ultimate:"name,criteria"`
	users []*User `ultimate:"users,aggregate"`
}
