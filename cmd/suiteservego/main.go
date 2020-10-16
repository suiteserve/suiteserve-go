package suiteservego

import (
	"flag"
	"net/http"
	"strings"
)

func main() {
	flag.Parse()

	http.Post("localhost:8080", "application/json", strings.NewReader(`
{

}
`))
}
