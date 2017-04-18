package pushaction

import "code.cloudfoundry.org/cli/actor/v2action"

func (actor Actor) DefaultDomain(orgGUID string) (v2action.Domain, Warnings, error) {
	return v2action.Domain{}, nil, nil
}
