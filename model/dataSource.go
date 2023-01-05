package model

type Begrep struct {
	ID   string `bson:"_id"`
	Term string `bson:"term"`
	Def  string `bson:"def"`
}
