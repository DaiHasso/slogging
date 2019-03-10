package logging

import (
    "github.com/pkg/errors"
)

type (
    // Extras defines a key, value map that will be marshaled to the final log
    // message.
    Extras map[string]interface{}

    // ValueFunc defines a function which will generate the eventual value that
    // will be logged for a ExtrasFunc map.
    ValueFunc func() (interface{}, error)

    // ExtrasGenerator is a function which produces an Extras map when called.
    // The resulting Extras (and possibly error) will be handled by the logger.
    ExtrasGenerator func() (Extras, error)

    // ExtrasFuncs is a map of key, values which will be evaluated at log-time
    // to get the resulting log value.
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

// StaticExtras adds all items in a key, value (Extras) map that will append:
//   key: value
// to the log for each item in the map.
// Refer to log formats specifications for how this will look when logged.
func StaticExtras(attributes Extras) ExtrasGenerator {
    return func() (Extras, error) {
        return attributes, nil
    }
}

// FunctionalExtras adds all items in a key, value (function) map that
// will append:
//   key: value()
// to the log for each item in the map.
// Refer to log formats specifications for how this will look when logged.
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
