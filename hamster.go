package hamster

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    "strings"

    "github.com/pkg/errors"
)

type Hamster struct {
    db     *sql.DB
    option *Option
}

type Option struct{}

type Food struct {
    Table    string
    Filepath string
}

func New(db *sql.DB, option *Option) *Hamster {
    return &Hamster{
        db:     db,
        option: option,
    }
}

func (ham *Hamster) Stuff(feed []*Food) error {
    tx, err := ham.db.Begin()
    if err != nil {
        return errors.Wrap(err, "error begin transaction")
    }

    for _, f := range feed {
        columns, rows, err := ham.loadCSV(f.Filepath)
        if err != nil {
            return errors.Wrap(err, "error load csv")
        }

        if err := ham.importData(f.Table, columns, rows); err != nil {
            return errors.Wrap(err, "error import data")
        }

        log.Printf("[DONE] imported %s\n", f.Table)
    }

    tx.Commit()
    return nil
}

func (ham *Hamster) importData(table string, columns []string, rows [][]string) error {
    colStrings := make([]string, len(columns))
    for i := range columns {
        colStrings[i] = "?"
    }
    colStringSet := fmt.Sprintf("(%s)", strings.Join(colStrings, ","))

    valStrings := make([]string, len(rows))
    values := make([]interface{}, 0, len(columns)*len(rows))
    for i, row := range rows {
        valStrings[i] = colStringSet

        for _, cell := range row {
            values = append(values, cell)
        }
    }

    if _, err := ham.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
        return errors.Wrap(err, "error exec query")
    }

    if len(valStrings) > 0 {
        stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, strings.Join(columns, ","), strings.Join(valStrings, ","))
        if _, err := ham.db.Exec(stmt, values...); err != nil {
            return errors.Wrap(err, "error exec query")
        }
    }
    return nil
}

func (ham *Hamster) loadCSV(filepath string) ([]string, [][]string, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return nil, nil, errors.Wrapf(err, "error os open file, file path: %s", filepath)
    }
    defer f.Close()

    reader := csv.NewReader(f)
    columns, err := reader.Read()
    if err != nil {
        return nil, nil, errors.Wrapf(err, "error read first line, file path: %s", filepath)
    }

    var rows [][]string
    for {
        row, err := reader.Read()
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, nil, errors.Wrapf(err, "error read line, file path: %s", filepath)
        }
        rows = append(rows, row)
    }
    return columns, rows, nil
}
