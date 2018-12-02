package slogging

import (
	"github.com/pkg/errors"
)

type ExtraParameter func(LogInstance) error

func AddExtraAttributes(attributes map[string]interface{}) ExtraParameter {
	return func(log LogInstance) error {
		for key, val := range attributes {
			log.With(key, val)
		}

		return nil
	}
}

func applyGlobalExtras(log LogInstance) error {
	globalExtrasMutex.RLock()
	defer globalExtrasMutex.RUnlock()
	for _, extra := range(globalExtras) {
		err := extra(log)
		if err != nil {
			return errors.Wrap(err, "Error while applying global extras")
		}
	}

	return nil
}
