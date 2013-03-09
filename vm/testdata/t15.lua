-- Test CONCAT with 5 values, one non-string float number
local a, b, c
a = "test"
b = "some"
c = a .. b .. "ugly" .. 123.45 .. 14
