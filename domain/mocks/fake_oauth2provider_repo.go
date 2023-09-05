// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
)

type FakeOauth2ProviderRepo struct {
	GetGoogleOauthTokenStub        func(string) (*model.GoogleOauth2Token, error)
	getGoogleOauthTokenMutex       sync.RWMutex
	getGoogleOauthTokenArgsForCall []struct {
		arg1 string
	}
	getGoogleOauthTokenReturns struct {
		result1 *model.GoogleOauth2Token
		result2 error
	}
	getGoogleOauthTokenReturnsOnCall map[int]struct {
		result1 *model.GoogleOauth2Token
		result2 error
	}
	GetGoogleOauthUserStub        func(*model.GoogleOauth2Token) (*model.GoogleOauth2User, error)
	getGoogleOauthUserMutex       sync.RWMutex
	getGoogleOauthUserArgsForCall []struct {
		arg1 *model.GoogleOauth2Token
	}
	getGoogleOauthUserReturns struct {
		result1 *model.GoogleOauth2User
		result2 error
	}
	getGoogleOauthUserReturnsOnCall map[int]struct {
		result1 *model.GoogleOauth2User
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthToken(arg1 string) (*model.GoogleOauth2Token, error) {
	fake.getGoogleOauthTokenMutex.Lock()
	ret, specificReturn := fake.getGoogleOauthTokenReturnsOnCall[len(fake.getGoogleOauthTokenArgsForCall)]
	fake.getGoogleOauthTokenArgsForCall = append(fake.getGoogleOauthTokenArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetGoogleOauthTokenStub
	fakeReturns := fake.getGoogleOauthTokenReturns
	fake.recordInvocation("GetGoogleOauthToken", []interface{}{arg1})
	fake.getGoogleOauthTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthTokenCallCount() int {
	fake.getGoogleOauthTokenMutex.RLock()
	defer fake.getGoogleOauthTokenMutex.RUnlock()
	return len(fake.getGoogleOauthTokenArgsForCall)
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthTokenCalls(stub func(string) (*model.GoogleOauth2Token, error)) {
	fake.getGoogleOauthTokenMutex.Lock()
	defer fake.getGoogleOauthTokenMutex.Unlock()
	fake.GetGoogleOauthTokenStub = stub
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthTokenArgsForCall(i int) string {
	fake.getGoogleOauthTokenMutex.RLock()
	defer fake.getGoogleOauthTokenMutex.RUnlock()
	argsForCall := fake.getGoogleOauthTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthTokenReturns(result1 *model.GoogleOauth2Token, result2 error) {
	fake.getGoogleOauthTokenMutex.Lock()
	defer fake.getGoogleOauthTokenMutex.Unlock()
	fake.GetGoogleOauthTokenStub = nil
	fake.getGoogleOauthTokenReturns = struct {
		result1 *model.GoogleOauth2Token
		result2 error
	}{result1, result2}
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthTokenReturnsOnCall(i int, result1 *model.GoogleOauth2Token, result2 error) {
	fake.getGoogleOauthTokenMutex.Lock()
	defer fake.getGoogleOauthTokenMutex.Unlock()
	fake.GetGoogleOauthTokenStub = nil
	if fake.getGoogleOauthTokenReturnsOnCall == nil {
		fake.getGoogleOauthTokenReturnsOnCall = make(map[int]struct {
			result1 *model.GoogleOauth2Token
			result2 error
		})
	}
	fake.getGoogleOauthTokenReturnsOnCall[i] = struct {
		result1 *model.GoogleOauth2Token
		result2 error
	}{result1, result2}
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUser(arg1 *model.GoogleOauth2Token) (*model.GoogleOauth2User, error) {
	fake.getGoogleOauthUserMutex.Lock()
	ret, specificReturn := fake.getGoogleOauthUserReturnsOnCall[len(fake.getGoogleOauthUserArgsForCall)]
	fake.getGoogleOauthUserArgsForCall = append(fake.getGoogleOauthUserArgsForCall, struct {
		arg1 *model.GoogleOauth2Token
	}{arg1})
	stub := fake.GetGoogleOauthUserStub
	fakeReturns := fake.getGoogleOauthUserReturns
	fake.recordInvocation("GetGoogleOauthUser", []interface{}{arg1})
	fake.getGoogleOauthUserMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUserCallCount() int {
	fake.getGoogleOauthUserMutex.RLock()
	defer fake.getGoogleOauthUserMutex.RUnlock()
	return len(fake.getGoogleOauthUserArgsForCall)
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUserCalls(stub func(*model.GoogleOauth2Token) (*model.GoogleOauth2User, error)) {
	fake.getGoogleOauthUserMutex.Lock()
	defer fake.getGoogleOauthUserMutex.Unlock()
	fake.GetGoogleOauthUserStub = stub
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUserArgsForCall(i int) *model.GoogleOauth2Token {
	fake.getGoogleOauthUserMutex.RLock()
	defer fake.getGoogleOauthUserMutex.RUnlock()
	argsForCall := fake.getGoogleOauthUserArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUserReturns(result1 *model.GoogleOauth2User, result2 error) {
	fake.getGoogleOauthUserMutex.Lock()
	defer fake.getGoogleOauthUserMutex.Unlock()
	fake.GetGoogleOauthUserStub = nil
	fake.getGoogleOauthUserReturns = struct {
		result1 *model.GoogleOauth2User
		result2 error
	}{result1, result2}
}

func (fake *FakeOauth2ProviderRepo) GetGoogleOauthUserReturnsOnCall(i int, result1 *model.GoogleOauth2User, result2 error) {
	fake.getGoogleOauthUserMutex.Lock()
	defer fake.getGoogleOauthUserMutex.Unlock()
	fake.GetGoogleOauthUserStub = nil
	if fake.getGoogleOauthUserReturnsOnCall == nil {
		fake.getGoogleOauthUserReturnsOnCall = make(map[int]struct {
			result1 *model.GoogleOauth2User
			result2 error
		})
	}
	fake.getGoogleOauthUserReturnsOnCall[i] = struct {
		result1 *model.GoogleOauth2User
		result2 error
	}{result1, result2}
}

func (fake *FakeOauth2ProviderRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getGoogleOauthTokenMutex.RLock()
	defer fake.getGoogleOauthTokenMutex.RUnlock()
	fake.getGoogleOauthUserMutex.RLock()
	defer fake.getGoogleOauthUserMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOauth2ProviderRepo) recordInvocation(key string, args []interface{}) {
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

var _ repository.Oauth2ProviderRepo = new(FakeOauth2ProviderRepo)