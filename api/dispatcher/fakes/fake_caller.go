package fakes

import (
	bocaction "github.com/frodenas/bosh-openstack-cpi/action"
)

type FakeCaller struct {
	CallAction bocaction.Action
	CallArgs   []interface{}
	CallResult interface{}
	CallErr    error
}

func (caller *FakeCaller) Call(action bocaction.Action, args []interface{}) (interface{}, error) {
	caller.CallAction = action
	caller.CallArgs = args
	return caller.CallResult, caller.CallErr
}
