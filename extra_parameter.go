package logging

import (
    "github.com/pkg/errors"
)

type (
    Extras map[string]interface{}
    ValueFunc func() (interface{}, error)
    ExtrasGenerator func() (Extras, error)
    ExtrasFuncs map[string]ValueFunc
)

// Extra is a convinience function for creating an Extras map it has some
// overhead of generating several extras that will have to be combined but it
// may be prefereable at times.
func Extra(key string, value interface{}) Extras {
    return map[string]interface{}{
        key: value,
    }
}

// AddExtraAttributes adds all items in a key, value map that will append:
//   key: value
// to the log for each item in the map.
func StaticExtras(attributes Extras) ExtrasGenerator {
    return func() (Extras, error) {
        return attributes, nil
    }
}

// AddExtraValueFuncs adds all items in a key, value (function) map that
// will append:
//   key: value()
// to the log for each item in the map.
func FunctionalExtras(
    attributeFuncs ExtrasFuncs,
) ExtrasGenerator {
    return func() (Extras, error) {
        extras := make(Extras)
        for key, val := range attributeFuncs {
            val, err := val()
            if err != nil {
                return extras, errors.Wrap(
                    err, "Error while evaluating attribute for log.",
                )
            }
            extras[key] = val
        }

        return extras, nil
    }
}

func applyExtrasToLogInstance(
    log LogInstance, extraParams []ExtrasGenerator,
) error {
    for _, extra := range(extraParams) {
        extras, err := extra()
        if err != nil {
            return errors.Wrap(
                err, "Error while applying global extras to log instance",
            )
        }
        for key, val := range extras {
            log.With(key, val)
        }
    }

    return nil
}

func applyGlobalExtrasToLogInstance(log LogInstance) error {
    globalExtraGeneratorsMutex.RLock()
    defer globalExtraGeneratorsMutex.RUnlock()
    return applyExtrasToLogInstance(log, globalExtraGenerators)
}
