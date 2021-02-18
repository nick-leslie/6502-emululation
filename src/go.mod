module local.com/main

go 1.15


replace(
local.com/memory => ./memory
local.com/cpu => ./cpu
)

require(
local.com/memory v0.0.0
local.com/cpu v0.0.0
)