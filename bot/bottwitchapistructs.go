package bot

type chattersData struct {
	Chatter_count int `json:"chatter_count"`
	Chatters      struct {
		Broadcaster []string `json:"broadcaster"`
		Vips        []string `json:"vips"`
		Moderators  []string `json:"moderators"`
		Staff       []string `json:"staff"`
		Admins      []string `json:"admins"`
		Global_mods []string `json:"global_mods"`
		Viewers     []string `json:"viewers"`
	} `json:"chatters"`
}

type usersData struct {
	User []struct {
		Id                string `json:"id"`
		Login             string `json:"login"`
		Display_name      string `json:"display_name"`
		Type              string `json:"type"`
		Broadcaster_type  string `json:"broadcaster_type"`
		Description       string `json:"description"`
		Profile_image_url string `json:"profile_image_url"`
		Offline_image_url string `json:"offline_image_url"`
		View_count        int    `json:"view_count"`
		Email             string `json:"email"`
	} `json:"data"`
}

type broadcasterSubscriptionsData struct {
	Subscriptions []struct {
		Broadcaster_id   string `json:"broadcaster_id"`
		Broadcaster_name string `json:"broadcaster_name"`
		Is_gift          bool   `json:"is_gift"`
		Tier             string `json:"tier"`
		Plan_name        string `json:"plan_name"`
		User_id          string `json:"user_id"`
		User_name        string `json:"user_name"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}
