package confs


type configuration struct {
	AbciDaemon     string
}

var Conf = configuration{}


func init() {
	Conf.AbciDaemon = "tcp://0.0.0.0:26658"
}
