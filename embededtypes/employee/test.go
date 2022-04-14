package employee

import "testing"

type fakeStmtForMaleCount struct {
	Stmt
}

func (f *fakeStmtForMaleCount) Exec(stmt string, args ...string) (employee.Result, error) {
	return Result{Count: 5}, nil
}

func TestEmployeeMaleCount(t, *testing.T) {
	f := fakeStmtForMaleCount{}
	c, _ := MaleCount(f)
	if c != 5 {
		t.Errorf("want: %d, actual: %s", 5, c)
		return
	}
}
