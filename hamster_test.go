package hamster

import (
    "reflect"
    "testing"
)

func Test_loadCSV(t *testing.T) {
    testFilePath := "./test/sample.csv"

    expectedColums := []string{
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
    resultColums, resultRows, err := ham.loadCSV(testFilePath)
    if err != nil {
        t.Fatal("error load csv", err.Error())
    }

    if !reflect.DeepEqual(resultColums, expectedColums) {
        t.Errorf("error colums %v, want %v", resultColums, expectedColums)
    }
    if !reflect.DeepEqual(resultRows, expectedRows) {
        t.Errorf("error rows %v, want %v", resultRows, expectedRows)
    }
}
