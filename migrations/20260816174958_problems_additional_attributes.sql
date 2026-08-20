-- +goose Up
CREATE TYPE complexity AS ENUM ('easy', 'medium', 'hard');

CREATE TYPE category AS ENUM ('array', 'hash map', 'two pointers', 'stack', 'binary search', 'sliding window', 'linked list', 'tree', 'trie', 'backtracking',
'heap', 'graph', '1d dynamic programming', 'bit manipulation', 'math', '2d dynamic programming', 'advanced graph', 'greedy', 'interval');

ALTER TABLE public.problems
ADD COLUMN complexity_rate complexity NOT NULL;

ALTER TABLE public.problems
ADD COLUMN categories category[] DEFAULT '{}' NOT NULL;

ALTER TABLE public.code_verification
ADD COLUMN code_template VARCHAR(1000);



-- +goose Down
ALTER TABLE public.code_verification
DROP COLUMN code_template;

ALTER TABLE public.problems
DROP COLUMN categories;

ALTER TABLE public.problems
DROP COLUMN complexity_rate;

DROP TYPE IF EXISTS category;

DROP TYPE IF EXISTS complexity;
