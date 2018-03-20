package app

// GCEZoneList is JSON obj from API `https://compute.googleapis.com/compute/v1/projects/PROJECT_ID/zones` resp body
type GCEZoneList struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Items []struct {
		Name string `json:"name"`
	} `json:"items"`
	SelfLink string `json:"selfLink"`
}

// GCEDisk is JSON obj from API `https://www.googleapis.com/compute/v1/projects/PROJECT_ID/zones/europe-west4-b/disks` resp body
type GCEDisk struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Items []struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	} `json:"items"`
	SelfLink string `json:"selfLink"`
}
