-- Test CALL of CLOSURE using Fibonacci algorithm with simple value of 2

-- Copyright (c) 2013 the authors listed at the following URL, and/or
-- the authors of referenced articles or incorporated external code:
-- http://en.literateprograms.org/Fibonacci_numbers_(Lua)?action=history&offset=20120305215844


function fib(n)
  return n<2 and n or fib(n-1)+fib(n-2)
end

a = fib(2)
