![image](https://user-images.githubusercontent.com/15713787/49335480-5550f700-f631-11e8-86eb-15464a3313cd.png)

# Hamster

Load CSV and insert to database

[see example](https://github.com/uenoryo/hamster/blob/master/example/main.go)



```go

ham := hamser.New(db, &hamster.Option{})
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

```
