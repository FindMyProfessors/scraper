package ucf

type Course struct {
	CoursePrefix  string `json:"prefix"`
	CatalogNumber string `json:"catalog_number"`
	CourseTitle   string `json:"title"`
	NameFirst     string `json:"name_first"`
	NameLast      string `json:"name_last"`
}
