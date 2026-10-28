package telegram

import "errors"

// ErrConfigInvalid reports that the bot's configuration cannot be used as
// supplied. The concrete problem travels as wrapped detail in the message;
// a detailed sentinel would only be justified if application logic branched
// on this specific failure.
var ErrConfigInvalid = errors.New("config invalid")
