package err

type ResultStatus int32

const (
	RsSuccess ResultStatus = iota
	RsFail
	RsNotValid
	RsDuplicateName
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
