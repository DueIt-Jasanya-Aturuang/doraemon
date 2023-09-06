// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"database/sql"
	"sync"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
)

type FakeAccessSqlRepo struct {
	CloseConnStub        func()
	closeConnMutex       sync.RWMutex
	closeConnArgsForCall []struct {
	}
	CreateAccessStub        func(context.Context, *model.Access) (*model.Access, error)
	createAccessMutex       sync.RWMutex
	createAccessArgsForCall []struct {
		arg1 context.Context
		arg2 *model.Access
	}
	createAccessReturns struct {
		result1 *model.Access
		result2 error
	}
	createAccessReturnsOnCall map[int]struct {
		result1 *model.Access
		result2 error
	}
	EndTxStub        func(error) error
	endTxMutex       sync.RWMutex
	endTxArgsForCall []struct {
		arg1 error
	}
	endTxReturns struct {
		result1 error
	}
	endTxReturnsOnCall map[int]struct {
		result1 error
	}
	GetAccessByUserIDAndAppIDStub        func(context.Context, string, string) (*model.Access, error)
	getAccessByUserIDAndAppIDMutex       sync.RWMutex
	getAccessByUserIDAndAppIDArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getAccessByUserIDAndAppIDReturns struct {
		result1 *model.Access
		result2 error
	}
	getAccessByUserIDAndAppIDReturnsOnCall map[int]struct {
		result1 *model.Access
		result2 error
	}
	GetConnStub        func() (*sql.Conn, error)
	getConnMutex       sync.RWMutex
	getConnArgsForCall []struct {
	}
	getConnReturns struct {
		result1 *sql.Conn
		result2 error
	}
	getConnReturnsOnCall map[int]struct {
		result1 *sql.Conn
		result2 error
	}
	GetTxStub        func() (*sql.Tx, error)
	getTxMutex       sync.RWMutex
	getTxArgsForCall []struct {
	}
	getTxReturns struct {
		result1 *sql.Tx
		result2 error
	}
	getTxReturnsOnCall map[int]struct {
		result1 *sql.Tx
		result2 error
	}
	OpenConnStub        func(context.Context) error
	openConnMutex       sync.RWMutex
	openConnArgsForCall []struct {
		arg1 context.Context
	}
	openConnReturns struct {
		result1 error
	}
	openConnReturnsOnCall map[int]struct {
		result1 error
	}
	StartTxStub        func(context.Context, *sql.TxOptions) error
	startTxMutex       sync.RWMutex
	startTxArgsForCall []struct {
		arg1 context.Context
		arg2 *sql.TxOptions
	}
	startTxReturns struct {
		result1 error
	}
	startTxReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAccessSqlRepo) CloseConn() {
	fake.closeConnMutex.Lock()
	fake.closeConnArgsForCall = append(fake.closeConnArgsForCall, struct {
	}{})
	stub := fake.CloseConnStub
	fake.recordInvocation("CloseConn", []interface{}{})
	fake.closeConnMutex.Unlock()
	if stub != nil {
		fake.CloseConnStub()
	}
}

func (fake *FakeAccessSqlRepo) CloseConnCallCount() int {
	fake.closeConnMutex.RLock()
	defer fake.closeConnMutex.RUnlock()
	return len(fake.closeConnArgsForCall)
}

func (fake *FakeAccessSqlRepo) CloseConnCalls(stub func()) {
	fake.closeConnMutex.Lock()
	defer fake.closeConnMutex.Unlock()
	fake.CloseConnStub = stub
}

func (fake *FakeAccessSqlRepo) CreateAccess(arg1 context.Context, arg2 *model.Access) (*model.Access, error) {
	fake.createAccessMutex.Lock()
	ret, specificReturn := fake.createAccessReturnsOnCall[len(fake.createAccessArgsForCall)]
	fake.createAccessArgsForCall = append(fake.createAccessArgsForCall, struct {
		arg1 context.Context
		arg2 *model.Access
	}{arg1, arg2})
	stub := fake.CreateAccessStub
	fakeReturns := fake.createAccessReturns
	fake.recordInvocation("CreateAccess", []interface{}{arg1, arg2})
	fake.createAccessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccessSqlRepo) CreateAccessCallCount() int {
	fake.createAccessMutex.RLock()
	defer fake.createAccessMutex.RUnlock()
	return len(fake.createAccessArgsForCall)
}

func (fake *FakeAccessSqlRepo) CreateAccessCalls(stub func(context.Context, *model.Access) (*model.Access, error)) {
	fake.createAccessMutex.Lock()
	defer fake.createAccessMutex.Unlock()
	fake.CreateAccessStub = stub
}

func (fake *FakeAccessSqlRepo) CreateAccessArgsForCall(i int) (context.Context, *model.Access) {
	fake.createAccessMutex.RLock()
	defer fake.createAccessMutex.RUnlock()
	argsForCall := fake.createAccessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAccessSqlRepo) CreateAccessReturns(result1 *model.Access, result2 error) {
	fake.createAccessMutex.Lock()
	defer fake.createAccessMutex.Unlock()
	fake.CreateAccessStub = nil
	fake.createAccessReturns = struct {
		result1 *model.Access
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) CreateAccessReturnsOnCall(i int, result1 *model.Access, result2 error) {
	fake.createAccessMutex.Lock()
	defer fake.createAccessMutex.Unlock()
	fake.CreateAccessStub = nil
	if fake.createAccessReturnsOnCall == nil {
		fake.createAccessReturnsOnCall = make(map[int]struct {
			result1 *model.Access
			result2 error
		})
	}
	fake.createAccessReturnsOnCall[i] = struct {
		result1 *model.Access
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) EndTx(arg1 error) error {
	fake.endTxMutex.Lock()
	ret, specificReturn := fake.endTxReturnsOnCall[len(fake.endTxArgsForCall)]
	fake.endTxArgsForCall = append(fake.endTxArgsForCall, struct {
		arg1 error
	}{arg1})
	stub := fake.EndTxStub
	fakeReturns := fake.endTxReturns
	fake.recordInvocation("EndTx", []interface{}{arg1})
	fake.endTxMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeAccessSqlRepo) EndTxCallCount() int {
	fake.endTxMutex.RLock()
	defer fake.endTxMutex.RUnlock()
	return len(fake.endTxArgsForCall)
}

func (fake *FakeAccessSqlRepo) EndTxCalls(stub func(error) error) {
	fake.endTxMutex.Lock()
	defer fake.endTxMutex.Unlock()
	fake.EndTxStub = stub
}

func (fake *FakeAccessSqlRepo) EndTxArgsForCall(i int) error {
	fake.endTxMutex.RLock()
	defer fake.endTxMutex.RUnlock()
	argsForCall := fake.endTxArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAccessSqlRepo) EndTxReturns(result1 error) {
	fake.endTxMutex.Lock()
	defer fake.endTxMutex.Unlock()
	fake.EndTxStub = nil
	fake.endTxReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) EndTxReturnsOnCall(i int, result1 error) {
	fake.endTxMutex.Lock()
	defer fake.endTxMutex.Unlock()
	fake.EndTxStub = nil
	if fake.endTxReturnsOnCall == nil {
		fake.endTxReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.endTxReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppID(arg1 context.Context, arg2 string, arg3 string) (*model.Access, error) {
	fake.getAccessByUserIDAndAppIDMutex.Lock()
	ret, specificReturn := fake.getAccessByUserIDAndAppIDReturnsOnCall[len(fake.getAccessByUserIDAndAppIDArgsForCall)]
	fake.getAccessByUserIDAndAppIDArgsForCall = append(fake.getAccessByUserIDAndAppIDArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetAccessByUserIDAndAppIDStub
	fakeReturns := fake.getAccessByUserIDAndAppIDReturns
	fake.recordInvocation("GetAccessByUserIDAndAppID", []interface{}{arg1, arg2, arg3})
	fake.getAccessByUserIDAndAppIDMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppIDCallCount() int {
	fake.getAccessByUserIDAndAppIDMutex.RLock()
	defer fake.getAccessByUserIDAndAppIDMutex.RUnlock()
	return len(fake.getAccessByUserIDAndAppIDArgsForCall)
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppIDCalls(stub func(context.Context, string, string) (*model.Access, error)) {
	fake.getAccessByUserIDAndAppIDMutex.Lock()
	defer fake.getAccessByUserIDAndAppIDMutex.Unlock()
	fake.GetAccessByUserIDAndAppIDStub = stub
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppIDArgsForCall(i int) (context.Context, string, string) {
	fake.getAccessByUserIDAndAppIDMutex.RLock()
	defer fake.getAccessByUserIDAndAppIDMutex.RUnlock()
	argsForCall := fake.getAccessByUserIDAndAppIDArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppIDReturns(result1 *model.Access, result2 error) {
	fake.getAccessByUserIDAndAppIDMutex.Lock()
	defer fake.getAccessByUserIDAndAppIDMutex.Unlock()
	fake.GetAccessByUserIDAndAppIDStub = nil
	fake.getAccessByUserIDAndAppIDReturns = struct {
		result1 *model.Access
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) GetAccessByUserIDAndAppIDReturnsOnCall(i int, result1 *model.Access, result2 error) {
	fake.getAccessByUserIDAndAppIDMutex.Lock()
	defer fake.getAccessByUserIDAndAppIDMutex.Unlock()
	fake.GetAccessByUserIDAndAppIDStub = nil
	if fake.getAccessByUserIDAndAppIDReturnsOnCall == nil {
		fake.getAccessByUserIDAndAppIDReturnsOnCall = make(map[int]struct {
			result1 *model.Access
			result2 error
		})
	}
	fake.getAccessByUserIDAndAppIDReturnsOnCall[i] = struct {
		result1 *model.Access
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) GetConn() (*sql.Conn, error) {
	fake.getConnMutex.Lock()
	ret, specificReturn := fake.getConnReturnsOnCall[len(fake.getConnArgsForCall)]
	fake.getConnArgsForCall = append(fake.getConnArgsForCall, struct {
	}{})
	stub := fake.GetConnStub
	fakeReturns := fake.getConnReturns
	fake.recordInvocation("GetConn", []interface{}{})
	fake.getConnMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccessSqlRepo) GetConnCallCount() int {
	fake.getConnMutex.RLock()
	defer fake.getConnMutex.RUnlock()
	return len(fake.getConnArgsForCall)
}

func (fake *FakeAccessSqlRepo) GetConnCalls(stub func() (*sql.Conn, error)) {
	fake.getConnMutex.Lock()
	defer fake.getConnMutex.Unlock()
	fake.GetConnStub = stub
}

func (fake *FakeAccessSqlRepo) GetConnReturns(result1 *sql.Conn, result2 error) {
	fake.getConnMutex.Lock()
	defer fake.getConnMutex.Unlock()
	fake.GetConnStub = nil
	fake.getConnReturns = struct {
		result1 *sql.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) GetConnReturnsOnCall(i int, result1 *sql.Conn, result2 error) {
	fake.getConnMutex.Lock()
	defer fake.getConnMutex.Unlock()
	fake.GetConnStub = nil
	if fake.getConnReturnsOnCall == nil {
		fake.getConnReturnsOnCall = make(map[int]struct {
			result1 *sql.Conn
			result2 error
		})
	}
	fake.getConnReturnsOnCall[i] = struct {
		result1 *sql.Conn
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) GetTx() (*sql.Tx, error) {
	fake.getTxMutex.Lock()
	ret, specificReturn := fake.getTxReturnsOnCall[len(fake.getTxArgsForCall)]
	fake.getTxArgsForCall = append(fake.getTxArgsForCall, struct {
	}{})
	stub := fake.GetTxStub
	fakeReturns := fake.getTxReturns
	fake.recordInvocation("GetTx", []interface{}{})
	fake.getTxMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccessSqlRepo) GetTxCallCount() int {
	fake.getTxMutex.RLock()
	defer fake.getTxMutex.RUnlock()
	return len(fake.getTxArgsForCall)
}

func (fake *FakeAccessSqlRepo) GetTxCalls(stub func() (*sql.Tx, error)) {
	fake.getTxMutex.Lock()
	defer fake.getTxMutex.Unlock()
	fake.GetTxStub = stub
}

func (fake *FakeAccessSqlRepo) GetTxReturns(result1 *sql.Tx, result2 error) {
	fake.getTxMutex.Lock()
	defer fake.getTxMutex.Unlock()
	fake.GetTxStub = nil
	fake.getTxReturns = struct {
		result1 *sql.Tx
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) GetTxReturnsOnCall(i int, result1 *sql.Tx, result2 error) {
	fake.getTxMutex.Lock()
	defer fake.getTxMutex.Unlock()
	fake.GetTxStub = nil
	if fake.getTxReturnsOnCall == nil {
		fake.getTxReturnsOnCall = make(map[int]struct {
			result1 *sql.Tx
			result2 error
		})
	}
	fake.getTxReturnsOnCall[i] = struct {
		result1 *sql.Tx
		result2 error
	}{result1, result2}
}

func (fake *FakeAccessSqlRepo) OpenConn(arg1 context.Context) error {
	fake.openConnMutex.Lock()
	ret, specificReturn := fake.openConnReturnsOnCall[len(fake.openConnArgsForCall)]
	fake.openConnArgsForCall = append(fake.openConnArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.OpenConnStub
	fakeReturns := fake.openConnReturns
	fake.recordInvocation("OpenConn", []interface{}{arg1})
	fake.openConnMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeAccessSqlRepo) OpenConnCallCount() int {
	fake.openConnMutex.RLock()
	defer fake.openConnMutex.RUnlock()
	return len(fake.openConnArgsForCall)
}

func (fake *FakeAccessSqlRepo) OpenConnCalls(stub func(context.Context) error) {
	fake.openConnMutex.Lock()
	defer fake.openConnMutex.Unlock()
	fake.OpenConnStub = stub
}

func (fake *FakeAccessSqlRepo) OpenConnArgsForCall(i int) context.Context {
	fake.openConnMutex.RLock()
	defer fake.openConnMutex.RUnlock()
	argsForCall := fake.openConnArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAccessSqlRepo) OpenConnReturns(result1 error) {
	fake.openConnMutex.Lock()
	defer fake.openConnMutex.Unlock()
	fake.OpenConnStub = nil
	fake.openConnReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) OpenConnReturnsOnCall(i int, result1 error) {
	fake.openConnMutex.Lock()
	defer fake.openConnMutex.Unlock()
	fake.OpenConnStub = nil
	if fake.openConnReturnsOnCall == nil {
		fake.openConnReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.openConnReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) StartTx(arg1 context.Context, arg2 *sql.TxOptions) error {
	fake.startTxMutex.Lock()
	ret, specificReturn := fake.startTxReturnsOnCall[len(fake.startTxArgsForCall)]
	fake.startTxArgsForCall = append(fake.startTxArgsForCall, struct {
		arg1 context.Context
		arg2 *sql.TxOptions
	}{arg1, arg2})
	stub := fake.StartTxStub
	fakeReturns := fake.startTxReturns
	fake.recordInvocation("StartTx", []interface{}{arg1, arg2})
	fake.startTxMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeAccessSqlRepo) StartTxCallCount() int {
	fake.startTxMutex.RLock()
	defer fake.startTxMutex.RUnlock()
	return len(fake.startTxArgsForCall)
}

func (fake *FakeAccessSqlRepo) StartTxCalls(stub func(context.Context, *sql.TxOptions) error) {
	fake.startTxMutex.Lock()
	defer fake.startTxMutex.Unlock()
	fake.StartTxStub = stub
}

func (fake *FakeAccessSqlRepo) StartTxArgsForCall(i int) (context.Context, *sql.TxOptions) {
	fake.startTxMutex.RLock()
	defer fake.startTxMutex.RUnlock()
	argsForCall := fake.startTxArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAccessSqlRepo) StartTxReturns(result1 error) {
	fake.startTxMutex.Lock()
	defer fake.startTxMutex.Unlock()
	fake.StartTxStub = nil
	fake.startTxReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) StartTxReturnsOnCall(i int, result1 error) {
	fake.startTxMutex.Lock()
	defer fake.startTxMutex.Unlock()
	fake.StartTxStub = nil
	if fake.startTxReturnsOnCall == nil {
		fake.startTxReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.startTxReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAccessSqlRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeConnMutex.RLock()
	defer fake.closeConnMutex.RUnlock()
	fake.createAccessMutex.RLock()
	defer fake.createAccessMutex.RUnlock()
	fake.endTxMutex.RLock()
	defer fake.endTxMutex.RUnlock()
	fake.getAccessByUserIDAndAppIDMutex.RLock()
	defer fake.getAccessByUserIDAndAppIDMutex.RUnlock()
	fake.getConnMutex.RLock()
	defer fake.getConnMutex.RUnlock()
	fake.getTxMutex.RLock()
	defer fake.getTxMutex.RUnlock()
	fake.openConnMutex.RLock()
	defer fake.openConnMutex.RUnlock()
	fake.startTxMutex.RLock()
	defer fake.startTxMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAccessSqlRepo) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ repository.AccessSqlRepo = new(FakeAccessSqlRepo)
