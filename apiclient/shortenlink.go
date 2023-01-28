package apiclient

type ShortenLink struct {
	Error     bool   `json:"error"`
	Exception string `json:"exception"`
	Domain    string `json:"path"`
	Path      string `json:"path"`
	LongLink  string `json:"long_link"`
}

type RedirectionData struct {
	Domain    string `json:"domain"`
	Referrer  string `json:"referrer"`
	Ip        string `json:"ip"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
}

func (apiClient PunyshortAPI) FollowRedirection(data RedirectionData) (ShortenLink, error) {
	proxyConfig := ShortenLink{}
	_, err := apiClient.RequestMap("POST", "/v1/follow", data, &proxyConfig)
	if err != nil {
		return ShortenLink{}, err
	}
	return proxyConfig, nil
}
