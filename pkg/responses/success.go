package responses

func CreateSuccessResponse[D any](data D) Response[D] {
	return Response[D]{Data: data, IsSuccess: true}
}
