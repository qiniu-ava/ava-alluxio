package interfaces

type CMDFlag interface {
	Validate() error
}
