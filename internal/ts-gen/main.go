// TODO this is not complete. Need to find a way to store all the API structs in a single file
package ts_gen

import (
	"fmt"

	m "github.com/Joshua-Hwang/habits2share/cmd/http"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	t := typescriptify.New()
	t.CreateInterface = false
	t.BackupDir=""



	err := t.ConvertToFile("./frontend/src/models.ts")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("OK")
}

