create table
  referee_trainer.question (
    id bigint primary key generated always as identity,
    text text not null,
    rule text not null,
    question_number integer not null
  );

create table
  referee_trainer.choice (
    id bigint primary key generated always as identity,
    question_id bigint references question (id),
    option text not null,
    is_answer boolean not null
  );
