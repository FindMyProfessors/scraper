package schools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scraper/models"
)

type UCFScraper struct {
	school *models.School
}

func (U *UCFScraper) StartRMSScraper() {

}

func (U *UCFScraper) StartSchoolScraper() {
	modalities := []string{"W", "V", "M", "RS", "P"}

	for _, mode := range modalities {
		var url = fmt.Sprintf("https://cdl.ucf.edu/wp-content/themes/cdl/lib/course-search-ajax.php?call=classes&term=1710&prefix=&catalog=&title=&instructor=&career=&college=&department=&mode=%s&_=1631687421296", mode)
		fmt.Println("url=" + url)
		go func(url string) {
			response, err := http.Get(url)
			if err != nil {
				fmt.Println("error occurred: ", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println("error occurred: ", err)
				}
			}(response.Body)
			body, err := io.ReadAll(response.Body)

			var data interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("error occurred: ", err)
			}
			fmt.Println("data=", data)
		}(url)
	}
}

func (U UCFScraper) GetSchool() models.School {
	if U.school != nil {
		return *U.school
	}
	U.StartRMSScraper()
	U.StartSchoolScraper()
	U.school = &models.School{}

	return *U.school
}
