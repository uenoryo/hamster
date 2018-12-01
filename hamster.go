package hamster

import (
    "database/sql"
    "encoding/csv"
    "io"
    "os"

    "github.com/pkg/errors"
)

type Hamster struct {
    db     *sql.DB
    option *Option
}

type Option struct{}

func New(db *sql.DB, option *Option) *Hamster {
    return &Hamster{
        db:     db,
        option: option,
    }
}

func (ham *Hamster) Stuff(filePath, table string) error {
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
