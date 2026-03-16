package usecases

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hashedPassword string) (bool, error)
}
