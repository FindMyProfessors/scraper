package rmp

import "github.com/FindMyProfessors/scraper/model"

type RMPProfessor struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	RMPId     string `json:"id,omitempty"`
	Reviews   struct {
		Edges []struct {
			Rating model.Review `json:"node"`
		} `json:"edges"`
	} `json:"ratings"`
}

type RMPResponseModel struct {
	Data struct {
		NewSearch struct {
			Teachers struct {
				Professors []struct {
					Professor RMPProfessor `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"teachers"`
		} `json:"newSearch"`
	} `json:"data"`
}
