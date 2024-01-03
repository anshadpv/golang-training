package main

type Database interface {
	Insert(id int, data string) error
	Update(id int, data string) error
	//Delete()
}

type sample interface {
}

var maap = make(map[int]string)

func main() {

	var db Database

	db = &msSQL{}
	db.Insert(1, "hello")
	db.Update(1, "hey")
	db.Insert(2, "hello world")
	db = &mySQL{}
	db.Insert(3, "ok")
	db.Update(4, "hoo")
	db.Insert(3, "no")

}
