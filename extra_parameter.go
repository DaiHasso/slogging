package slogging

import (
	"github.com/pkg/errors"
)

type ExtraParameter func(LogInstance) error

// AddExtraAttributes adds all items in a key, value map that will append:
//   key: value
// to the log for each item in the map.
func AddExtraAttributes(attributes map[string]interface{}) ExtraParameter {
	return func(log LogInstance) error {
		for key, val := range attributes {
			log.With(key, val)
		}

		return nil
	}
}

// AddExtraAttributeFuncs adds all items in a key, value (function) map that
// will append:
//   key: value()
// to the log for each item in the map.
func AddExtraAttributeFuncs(
	attributeFuncs map[string]func() (interface{}, error),
) ExtraParameter {
	return func(log LogInstance) error {
		for key, val := range attributeFuncs {
			val, err := val()
			if err != nil {
				return errors.Wrap(
					err, "Error while evaluating attribute for log.",
				)
			}
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
