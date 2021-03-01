// put any functions that you want to be available to the interface (DatabaseRepo)

package dbrepo

func (m *postgresDBRepo) AllUsers() bool {
	return true
}
