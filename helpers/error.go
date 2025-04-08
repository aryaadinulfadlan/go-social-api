package helpers

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
