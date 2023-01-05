-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS expense_id_seq;

-- Table Definition
CREATE TABLE "expenses" (
    "id" int4 NOT NULL DEFAULT nextval('expense_id_seq'::regclass),
    "title" text,
    "amount" FLOAT,
	"note" TEXT,
	"tags" TEXT[]
    PRIMARY KEY ("id")
);

INSERT INTO "expenses" ("id", "title", "amount","note","tags") 
VALUES (1, 'test-title', 'test-content', 199,"test-note",ARRAY["test1","test2"]);