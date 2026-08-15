-- +goose Up
CREATE TYPE available_languages AS ENUM ('go', 'py', 'java', 'cpp', 'js', 'ts', 'csharp', 'ko', 'sw', 'rust');

ALTER TABLE public.problems 
ADD COLUMN langs available_languages[] DEFAULT '{}';

ALTER TABLE public.code_verification
ADD COLUMN lang available_languages NOT NULL;

ALTER TABLE code_verification
ADD CONSTRAINT fk_problem_id
FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_languages_on_tests_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE public.problems 
    SET langs = array_append(COALESCE(langs, '{}'), NEW.lang) 
    WHERE id = NEW.problem_id 
      AND NOT (COALESCE(langs, '{}') @> ARRAY[NEW.lang]::available_languages[]);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_update_languages_on_tests_insert
AFTER INSERT ON public.code_verification
FOR EACH ROW
EXECUTE FUNCTION update_languages_on_tests_insert();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_languages_on_tests_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE public.problems 
    SET langs = array_remove(langs, OLD.lang) 
    WHERE id = OLD.problem_id; 
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_update_languages_on_tests_delete
AFTER DELETE ON public.code_verification
FOR EACH ROW
EXECUTE FUNCTION update_languages_on_tests_delete();

ALTER TABLE code_verification
ADD CONSTRAINT unique_problem_id_lang UNIQUE (problem_id, lang);

-- +goose Down
ALTER TABLE code_verification
DROP CONSTRAINT unique_problem_id_lang;

DROP TRIGGER IF EXISTS trigger_update_languages_on_tests_delete ON public.code_verification;
DROP FUNCTION update_languages_on_tests_delete();

DROP TRIGGER IF EXISTS trigger_update_languages_on_tests_insert ON public.code_verification;
DROP FUNCTION update_languages_on_tests_insert();

ALTER TABLE public.code_verification
DROP CONSTRAINT fk_problem_id;

ALTER TABLE public.code_verification
DROP COLUMN lang; 

ALTER TABLE public.problems
DROP COLUMN langs;

DROP TYPE IF EXISTS available_languages;
