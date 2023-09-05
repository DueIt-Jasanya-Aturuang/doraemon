// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
)

type FakeAccountApiRepo struct {
	CreateProfileStub        func([]byte) (*model.Profile, error)
	createProfileMutex       sync.RWMutex
	createProfileArgsForCall []struct {
		arg1 []byte
	}
	createProfileReturns struct {
		result1 *model.Profile
		result2 error
	}
	createProfileReturnsOnCall map[int]struct {
		result1 *model.Profile
		result2 error
	}
	GetProfileByUserIDStub        func(string) (*model.Profile, error)
	getProfileByUserIDMutex       sync.RWMutex
	getProfileByUserIDArgsForCall []struct {
		arg1 string
	}
	getProfileByUserIDReturns struct {
		result1 *model.Profile
		result2 error
	}
	getProfileByUserIDReturnsOnCall map[int]struct {
		result1 *model.Profile
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAccountApiRepo) CreateProfile(arg1 []byte) (*model.Profile, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.createProfileMutex.Lock()
	ret, specificReturn := fake.createProfileReturnsOnCall[len(fake.createProfileArgsForCall)]
	fake.createProfileArgsForCall = append(fake.createProfileArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	stub := fake.CreateProfileStub
	fakeReturns := fake.createProfileReturns
	fake.recordInvocation("CreateProfile", []interface{}{arg1Copy})
	fake.createProfileMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccountApiRepo) CreateProfileCallCount() int {
	fake.createProfileMutex.RLock()
	defer fake.createProfileMutex.RUnlock()
	return len(fake.createProfileArgsForCall)
}

func (fake *FakeAccountApiRepo) CreateProfileCalls(stub func([]byte) (*model.Profile, error)) {
	fake.createProfileMutex.Lock()
	defer fake.createProfileMutex.Unlock()
	fake.CreateProfileStub = stub
}

func (fake *FakeAccountApiRepo) CreateProfileArgsForCall(i int) []byte {
	fake.createProfileMutex.RLock()
	defer fake.createProfileMutex.RUnlock()
	argsForCall := fake.createProfileArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAccountApiRepo) CreateProfileReturns(result1 *model.Profile, result2 error) {
	fake.createProfileMutex.Lock()
	defer fake.createProfileMutex.Unlock()
	fake.CreateProfileStub = nil
	fake.createProfileReturns = struct {
		result1 *model.Profile
		result2 error
	}{result1, result2}
}

func (fake *FakeAccountApiRepo) CreateProfileReturnsOnCall(i int, result1 *model.Profile, result2 error) {
	fake.createProfileMutex.Lock()
	defer fake.createProfileMutex.Unlock()
	fake.CreateProfileStub = nil
	if fake.createProfileReturnsOnCall == nil {
		fake.createProfileReturnsOnCall = make(map[int]struct {
			result1 *model.Profile
			result2 error
		})
	}
	fake.createProfileReturnsOnCall[i] = struct {
		result1 *model.Profile
		result2 error
	}{result1, result2}
}

func (fake *FakeAccountApiRepo) GetProfileByUserID(arg1 string) (*model.Profile, error) {
	fake.getProfileByUserIDMutex.Lock()
	ret, specificReturn := fake.getProfileByUserIDReturnsOnCall[len(fake.getProfileByUserIDArgsForCall)]
	fake.getProfileByUserIDArgsForCall = append(fake.getProfileByUserIDArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetProfileByUserIDStub
	fakeReturns := fake.getProfileByUserIDReturns
	fake.recordInvocation("GetProfileByUserID", []interface{}{arg1})
	fake.getProfileByUserIDMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAccountApiRepo) GetProfileByUserIDCallCount() int {
	fake.getProfileByUserIDMutex.RLock()
	defer fake.getProfileByUserIDMutex.RUnlock()
	return len(fake.getProfileByUserIDArgsForCall)
}

func (fake *FakeAccountApiRepo) GetProfileByUserIDCalls(stub func(string) (*model.Profile, error)) {
	fake.getProfileByUserIDMutex.Lock()
	defer fake.getProfileByUserIDMutex.Unlock()
	fake.GetProfileByUserIDStub = stub
}

func (fake *FakeAccountApiRepo) GetProfileByUserIDArgsForCall(i int) string {
	fake.getProfileByUserIDMutex.RLock()
	defer fake.getProfileByUserIDMutex.RUnlock()
	argsForCall := fake.getProfileByUserIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAccountApiRepo) GetProfileByUserIDReturns(result1 *model.Profile, result2 error) {
	fake.getProfileByUserIDMutex.Lock()
	defer fake.getProfileByUserIDMutex.Unlock()
	fake.GetProfileByUserIDStub = nil
	fake.getProfileByUserIDReturns = struct {
		result1 *model.Profile
		result2 error
	}{result1, result2}
}

func (fake *FakeAccountApiRepo) GetProfileByUserIDReturnsOnCall(i int, result1 *model.Profile, result2 error) {
	fake.getProfileByUserIDMutex.Lock()
	defer fake.getProfileByUserIDMutex.Unlock()
	fake.GetProfileByUserIDStub = nil
	if fake.getProfileByUserIDReturnsOnCall == nil {
		fake.getProfileByUserIDReturnsOnCall = make(map[int]struct {
			result1 *model.Profile
			result2 error
		})
	}
	fake.getProfileByUserIDReturnsOnCall[i] = struct {
		result1 *model.Profile
		result2 error
	}{result1, result2}
}

func (fake *FakeAccountApiRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createProfileMutex.RLock()
	defer fake.createProfileMutex.RUnlock()
	fake.getProfileByUserIDMutex.RLock()
	defer fake.getProfileByUserIDMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAccountApiRepo) recordInvocation(key string, args []interface{}) {
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

var _ repository.AccountApiRepo = new(FakeAccountApiRepo)