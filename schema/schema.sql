create table
  question (
    id bigint primary key generated by default as identity,
    text text not null,
    rule text not null,
    question_number integer not null
  );

create table
  choice (
    id bigint primary key generated by default as identity,
    question_id bigint references question (id),
    option text not null,
    text text not null,
    is_answer boolean not null
  );
