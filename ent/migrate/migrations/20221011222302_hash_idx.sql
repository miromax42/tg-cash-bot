-- change index, because hash is lighter, faster and
-- the only comparisons supposed to make is eq

-- drop index "expense_created_by" from table: "expenses"
DROP INDEX "expense_created_by";
-- create index "expense_created_by" to table: "expenses"
CREATE INDEX "expense_created_by" ON "expenses" USING HASH ("created_by");
