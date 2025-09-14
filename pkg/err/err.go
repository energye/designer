package err

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
