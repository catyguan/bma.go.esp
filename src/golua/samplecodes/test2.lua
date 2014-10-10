-- local example = "an example string"
-- for i in string.gmatch(example, "%S+") do
--   print(i)
-- end

-- function deepcopy(orig)
--     local orig_type = type(orig)
--     local copy
--     if orig_type == 'table' then
--         copy = {}
--         for orig_key, orig_value in next, orig, nil do
--             copy[deepcopy(orig_key)] = deepcopy(orig_value)
--         end
--         setmetatable(copy, deepcopy(getmetatable(orig)))
--     else -- number, string, boolean, etc
--         copy = orig
--     end
--     return copy
-- end

-- function qsort(vec, low, high)
--   if low < high then
--     local middle = partition(vec, low, high)
--     qsort(vec, low, middle-1)
--     qsort(vec, middle+1, high)
--   end
-- end

-- Account = {}
-- Account.__index = Account

-- function Account.create(balance)
--    closure(Account)
--    local acnt = {}             -- our new object

--    setmetatable(acnt,Account)  -- make Account handle lookup
--    acnt.balance = balance      -- initialize our object
--    return acnt
-- end

-- function Account:withdraw(amount)
--    self.balance = self.balance - amount
-- end

-- -- create and use an Account
-- acc = Account.create(1000)
-- acc:withdraw(100)

rawset(Directions, "LEFT", 5)
print(Directions.LEFT)         -- prints 5
table.insert(Directions, 6)
print(Directions[1])           -- prints 6