package models

type Connection struct {
	Host      string
	Token     string
	Namespace string
	Insecure  bool
}
