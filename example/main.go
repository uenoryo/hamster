package main

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "github.com/uenoryo/hamster"
)

func main() {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hamster?charset=utf8")
    if err != nil {
        log.Fatalf(err.Error())
    }
    defer db.Close()

    ham := hamster.New(db, &hamster.Option{})
    feed := []*hamster.Food{
        {
            Table:    "hamsters",
            Filepath: "./data/sample01.csv",
        },
        {
            Table:    "feed",
            Filepath: "./data/sample02.csv",
        },
    }

    if err := ham.Stuff(feed); err != nil {
        log.Fatal(err.Error())
    }
}
