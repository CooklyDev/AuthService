package application

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hashedPassword string) (bool, error)
}

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	UserRepository() UserRepo
	SessionRepository() SessionRepo
}
