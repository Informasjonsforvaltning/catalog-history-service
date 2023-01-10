package model

type Concept struct {
	ID   int    `bson:"_id"`
	Term string `bson:"term"`
	Def  string `bson:"def"`
}
