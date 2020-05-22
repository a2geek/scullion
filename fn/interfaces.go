package fn

import "scullion/ctx"

// Registrar is a registration function to add functions dynamically into a State at runtime.
// All functions must return a value. If they are _methods_, then they should return error.
// The execution will look for an error return code, log the message, and stop that particular
// execution as appropriate.
type Registrar func(*ctx.State)
