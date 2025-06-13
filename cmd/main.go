package main

import (
	"fmt"

	"github.com/Pleiades-IUST/backend/utils/config"
	"github.com/Pleiades-IUST/backend/webservice"
)

func main() {

	r := webservice.SetupRouter()

	if err := r.Run(fmt.Sprintf("%s:%s", config.GetHost(), config.GetPort())); err != nil {
		panic(err)
	}

}
