module rootMod

go 1.16

replace(
local.com/memory => ./src/memory
local.com/cpu => ./src/cpu
)

require(
local.com/memory v0.0.0
local.com/cpu v0.0.0
)