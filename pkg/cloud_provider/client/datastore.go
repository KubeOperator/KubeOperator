package client

type DatastoreResult struct {
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	FreeSpace int    `json:"freeSpace"`
}
