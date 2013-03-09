-- Test SELF
o = {}
function o:add(n)
  return n + 1
end

a = o:add(5)
