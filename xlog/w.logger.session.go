package xlog

//CreateSession create logger session
func CreateSession() string {
	return uuid.New()
}
