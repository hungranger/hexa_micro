package shortener

type Redirect struct {
	Code     string `json:"code" bson:"code" msgpack:"code"`
	URL      string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreateAt int64  `json:"create_at" bson:"create_at" msgpack:"create_at"`
}
