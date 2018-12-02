package hamster

import (
    "database/sql"
    "fmt"
    "log"
    "reflect"
    "testing"

    _ "github.com/go-sql-driver/mysql"
)

func TestStaff(t *testing.T) {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hamster?charset=utf8")
    if err != nil {
        t.Fatal("error sql open", err.Error())
    }
    defer db.Close()

    testTables := []string{"hamsters", "feed"}

    type Test struct {
        Input          []*Food
        ExpectedCounts []int
        InitFunc       func()
    }

    tests := []Test{
        {
            Input: []*Food{
                {
                    Table:    testTables[0],
                    Filepath: "./test/sample01.csv",
                },
                {
                    Table:    testTables[1],
                    Filepath: "./test/sample02.csv",
                },
            },
            ExpectedCounts: []int{3, 5},
            InitFunc: func() {
                if _, err := db.Exec(fmt.Sprintf("CREATE TABLE %s(id int, name varchar(20), kind varchar(20))", testTables[0])); err != nil {
                    t.Fatal("error exec query", err.Error())
                }
                if _, err := db.Exec(fmt.Sprintf("CREATE TABLE %s(id int, name varchar(20), price int)", testTables[1])); err != nil {
                    t.Fatal("error exec query", err.Error())
                }
            },
        },
    }

    for i, test := range tests {
        t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
            defer func() {
                for _, ipt := range test.Input {
                    if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", ipt.Table)); err != nil {
                        t.Fatal("error exec query", err.Error())
                    }
                }
            }()

            test.InitFunc()

            ham := New(db, &Option{})
            if err := ham.Stuff(test.Input); err != nil {
                t.Fatal("error stuff", err.Error())
            }

            for i, ipt := range test.Input {
                rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM %s", ipt.Table))
                if err != nil {
                    t.Fatal("error exec query", err.Error())
                }
                defer rows.Close()

                var count int
                for rows.Next() {
                    if err := rows.Scan(&count); err != nil {
                        t.Fatal("error scan", err.Error())
                    }
                }

                if g, w := count, test.ExpectedCounts[i]; g != w {
                    t.Errorf("error table %s data count %d, want %d", ipt.Table, g, w)
                }
            }
        })
    }
}

func Test_importData(t *testing.T) {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hamster?charset=utf8")
    if err != nil {
        t.Fatal("error sql open", err.Error())
    }
    defer db.Close()

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
        if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", inputTable)); err != nil {
            t.Fatal("error exec query", err.Error())
        }
    }()

    if _, err := db.Exec(fmt.Sprintf("CREATE TABLE %s(id int, name varchar(20), kind varchar(20))", inputTable)); err != nil {
        t.Fatal("error exec query", err.Error())
    }

    ham := New(db, &Option{})
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
    testFilePath := "./test/sample01.csv"

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
