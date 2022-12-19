package model

type DataSource struct {
	Id                string `json:"id" bson:"_id"`
	DataSourceType    string `json:"dataSourceType" bson:"dataSourceType"`
	DataType          string `json:"dataType" bson:"dataType"`
	Url               string `json:"url" bson:"url"`
	AcceptHeaderValue string `json:"acceptHeaderValue" bson:"acceptHeaderValue"`
	PublisherId       string `json:"publisherId" bson:"publisherId"`
	Description       string `json:"description" bson:"description"`
}
