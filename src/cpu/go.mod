module local.com/cpu

go 1.15

replace(
local.com/memory => ../memory
)

require(
local.com/memory v0.0.0
)