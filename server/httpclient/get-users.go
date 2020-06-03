package httpclient

import "strconv"

// UserResult .
type UserResult struct {
	Data Datum `json:"data"`
	Ad   Ad    `json:"ad"`
}

// UsersResult .
type UsersResult struct {
	Page       int64   `json:"page"`
	PerPage    int64   `json:"per_page"`
	Total      int64   `json:"total"`
	TotalPages int64   `json:"total_pages"`
	Data       []Datum `json:"data"`
	Ad         Ad      `json:"ad"`
}

// Ad .
type Ad struct {
	Company string `json:"company"`
	URL     string `json:"url"`
	Text    string `json:"text"`
}

// Datum .
type Datum struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// GetUsers Get dummy users API
// https://reqres.in/api/users?page=1
func (c HTTPClient) GetUsers(page int) (*UsersResult, error) {
	result := &UsersResult{}
	_, err := c.R().
		SetResult(result).
		SetQueryParams(map[string]string{"page": strconv.Itoa(page)}).
		Get("users")

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUser Get single dummy user API
// https://reqres.in/api/users/2
func (c HTTPClient) GetUser(ID int) (*UserResult, error) {
	result := &UserResult{}
	_, err := c.R().
		SetResult(result).
		SetPathParams(map[string]string{"id": strconv.Itoa(ID)}).
		Get("users/{id}")

	if err != nil {
		return nil, err
	}

	return result, nil
}
