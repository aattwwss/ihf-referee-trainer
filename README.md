# ihf-referee-rules
All stuff regarding IHF rules and regulations.

All materials taken from the Internation Handball Federation (IHF) offical website [here](https://www.ihf.info).

[Rules](https://www.ihf.info/sites/default/files/2022-09/09A%20-%20Rules%20of%20the%20Game_Indoor%20Handball_E.pdf)

[Rules Questions](http://images.ihfeducation.ihf.info/File/Get?id=\ContentItems\Files\a\4\a4cb4e7a-c908-4dec-92fb-0348e4bbc05e.pdf)

[Rules Answers](http://images.ihfeducation.ihf.info/File/Get?id=\ContentItems\Files\e\d\edc086d8-0f70-45f4-bf40-1611ed8d8b50.pdf)

# Parsing the questions and answers PDF

I found the easiest way to handle the pdf is to simply render it in your browser then copy and paste the entire content into two separte text files, questions.txt and answers.txt.

With the two text files, we can attempt to parse them into workable formats, currently supporting `csv` and `json`.

```shell
# generate 2 csv files of questions and answers repectively
go run cmd/parse.go --f=csv

# generate a json array of all the questions, options and theirs answers
go run cmd/parse.go --f=csv
```
json output example
```json
[
    {
      "ID": 101,
      "Text": "BLACK 10 receives the ball at the free-throw line. Just after catching the ball he is pulled down by WHITE 3. He falls, hurts his elbow, and needs medical treatment on the court. Correct decision?",
      "Choices": [
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "a",
          "Text": "Warning for WHITE 3",
          "IsAnswer": false
        },
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "b",
          "Text": "2-minute suspension for WHITE 3",
          "IsAnswer": false
        },
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "c",
          "Text": "Two people from BLACK team, who are entitled to participate, can enter the court to give BLACK 10 medical treatment on the court after the hand signals 15 and 16 have been shown by one of the referees.",
          "IsAnswer": false
        },
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "d",
          "Text": "BLACK 10 may continue to play after receiving medical treatment on the court.",
          "IsAnswer": false
        },
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "e",
          "Text": "After receiving medical treatment on the court, BLACK 10 can only re-enter the court following the third attack of his team.",
          "IsAnswer": false
        },
        {
          "ID": 0,
          "QuestionID": 101,
          "Option": "f",
          "Text": "Time-out",
          "IsAnswer": false
        }
      ],
      "Rule": "274",
      "QuestionNum": 45
    }
]
```
