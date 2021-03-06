// This file was generated by counterfeiter
package appinstancesfakes

import (
	"sync"

	"code.cloudfoundry.org/cli/cf/api/appinstances"
	"code.cloudfoundry.org/cli/cf/models"
)

type FakeRepository struct {
	GetInstancesStub        func(appGUID string) (instances []models.AppInstanceFields, apiErr error)
	getInstancesMutex       sync.RWMutex
	getInstancesArgsForCall []struct {
		appGUID string
	}
	getInstancesReturns struct {
		result1 []models.AppInstanceFields
		result2 error
	}
	DeleteInstanceStub        func(appGUID string, instance int) error
	deleteInstanceMutex       sync.RWMutex
	deleteInstanceArgsForCall []struct {
		appGUID  string
		instance int
	}
	deleteInstanceReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRepository) GetInstances(appGUID string) (instances []models.AppInstanceFields, apiErr error) {
	fake.getInstancesMutex.Lock()
	fake.getInstancesArgsForCall = append(fake.getInstancesArgsForCall, struct {
		appGUID string
	}{appGUID})
	fake.recordInvocation("GetInstances", []interface{}{appGUID})
	fake.getInstancesMutex.Unlock()
	if fake.GetInstancesStub != nil {
		return fake.GetInstancesStub(appGUID)
	} else {
		return fake.getInstancesReturns.result1, fake.getInstancesReturns.result2
	}
}

func (fake *FakeRepository) GetInstancesCallCount() int {
	fake.getInstancesMutex.RLock()
	defer fake.getInstancesMutex.RUnlock()
	return len(fake.getInstancesArgsForCall)
}

func (fake *FakeRepository) GetInstancesArgsForCall(i int) string {
	fake.getInstancesMutex.RLock()
	defer fake.getInstancesMutex.RUnlock()
	return fake.getInstancesArgsForCall[i].appGUID
}

func (fake *FakeRepository) GetInstancesReturns(result1 []models.AppInstanceFields, result2 error) {
	fake.GetInstancesStub = nil
	fake.getInstancesReturns = struct {
		result1 []models.AppInstanceFields
		result2 error
	}{result1, result2}
}

func (fake *FakeRepository) DeleteInstance(appGUID string, instance int) error {
	fake.deleteInstanceMutex.Lock()
	fake.deleteInstanceArgsForCall = append(fake.deleteInstanceArgsForCall, struct {
		appGUID  string
		instance int
	}{appGUID, instance})
	fake.recordInvocation("DeleteInstance", []interface{}{appGUID, instance})
	fake.deleteInstanceMutex.Unlock()
	if fake.DeleteInstanceStub != nil {
		return fake.DeleteInstanceStub(appGUID, instance)
	} else {
		return fake.deleteInstanceReturns.result1
	}
}

func (fake *FakeRepository) DeleteInstanceCallCount() int {
	fake.deleteInstanceMutex.RLock()
	defer fake.deleteInstanceMutex.RUnlock()
	return len(fake.deleteInstanceArgsForCall)
}

func (fake *FakeRepository) DeleteInstanceArgsForCall(i int) (string, int) {
	fake.deleteInstanceMutex.RLock()
	defer fake.deleteInstanceMutex.RUnlock()
	return fake.deleteInstanceArgsForCall[i].appGUID, fake.deleteInstanceArgsForCall[i].instance
}

func (fake *FakeRepository) DeleteInstanceReturns(result1 error) {
	fake.DeleteInstanceStub = nil
	fake.deleteInstanceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getInstancesMutex.RLock()
	defer fake.getInstancesMutex.RUnlock()
	fake.deleteInstanceMutex.RLock()
	defer fake.deleteInstanceMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRepository) recordInvocation(key string, args []interface{}) {
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

var _ appinstances.Repository = new(FakeRepository)
