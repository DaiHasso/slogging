// Code generated by "stringer -type=LogTarget"; DO NOT EDIT.

package logging

import "strconv"

const _LogTarget_name = "StdoutLogFile"

var _LogTarget_index = [...]uint8{0, 6, 13}

func (i LogTarget) String() string {
	i -= 1
	if i < 0 || i >= LogTarget(len(_LogTarget_index)-1) {
		return "LogTarget(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _LogTarget_name[_LogTarget_index[i]:_LogTarget_index[i+1]]
}
