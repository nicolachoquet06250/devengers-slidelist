package routes

import (
	"html/template"
	"net/http"
)

type TestBluetoothApiPageParams struct {
}

func testBluetoothApiPage(w http.ResponseWriter, r *http.Request) {
	testBluetoothApiTpl, _ := TestBluetoothApiTpl()
	tpl, _ := template.New("auth").Parse(testBluetoothApiTpl)
	err := tpl.Execute(w, TestBluetoothApiPageParams{})

	if err != nil {
		println(err.Error())
	}
}
