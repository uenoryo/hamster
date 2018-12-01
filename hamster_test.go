package hamster

import (
    "database/sql"
    "fmt"
    "log"
    "reflect"
    "testing"

    _ "github.com/go-sql-driver/mysql"
)

func Test_importData(t *testing.T) {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hamster?charset=utf8")
    if err != nil {
        t.Fatal("error sql open", err.Error())
    }
    defer db.Close()

    ham := New(db, &Option{})

    inputColumns := []string{
        "id",
        "name",
        "kind",
    }
    inputRows := [][]string{
        {
            "1",
            "TOMATO",
            "Djungarian",
        },
        {
            "2",
            "DAIKON",
            "Pearl White",
        },
        {
            "3",
            "KAZUNOKO",
            "Golden",
        },
    }
    inputTable := "hamsters"

    defer func() {
        if _, err := ham.db.Exec(fmt.Sprintf("DROP TABLE %s", inputTable)); err != nil {
            t.Fatal("error exec query", err.Error())
        }
    }()

    if _, err := ham.db.Exec(fmt.Sprintf("CREATE TABLE %s(id int, name varchar(20), kind varchar(20))", inputTable)); err != nil {
        t.Fatal("error exec query", err.Error())
    }

    if err := ham.importData(inputTable, inputColumns, inputRows); err != nil {
        t.Fatal("error import data", err.Error())
    }

    rows, err := db.Query("SELECT id, name, kind FROM hamsters")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    i := 0
    for rows.Next() {
        var id, name, kind string
        if err := rows.Scan(&id, &name, &kind); err != nil {
            log.Fatal(err)
        }

        if g, w := id, inputRows[i][0]; g != w {
            t.Errorf("error row:%d id %s, want %s", i, g, w)
        }
        if g, w := name, inputRows[i][1]; g != w {
            t.Errorf("error row:%d name %s, want %s", i, g, w)
        }
        if g, w := kind, inputRows[i][2]; g != w {
            t.Errorf("error row:%d kind %s, want %s", i, g, w)
        }
        i++
    }

    if err := rows.Err(); err != nil {
        t.Fatal(err.Error())
    }
}

func Test_loadCSV(t *testing.T) {
    testFilePath := "./test/sample.csv"

    expectedColumns := []string{
        "id",
        "name",
        "kind",
    }
    expectedRows := [][]string{
        {
            "1",
            "TOMATO",
            "Djungarian",
        },
        {
            "2",
            "DAIKON",
            "Pearl White",
        },
        {
            "3",
            "KAZUNOKO",
            "Golden",
        },
    }

    ham := &Hamster{}
    resultColumns, resultRows, err := ham.loadCSV(testFilePath)
    if err != nil {
        t.Fatal("error load csv", err.Error())
    }

    if !reflect.DeepEqual(resultColumns, expectedColumns) {
        t.Errorf("error columns %v, want %v", resultColumns, expectedColumns)
    }
    if !reflect.DeepEqual(resultRows, expectedRows) {
        t.Errorf("error rows %v, want %v", resultRows, expectedRows)
    }
}
