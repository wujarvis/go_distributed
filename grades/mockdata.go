package grades

func init() {
	students = []Student{
		{
			ID: 1,
			FirstName: "Nick",
			LastName: "Carter",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type: GradeQuiz,
					Score: 83,
				},
				{
					Title: "First test",
					Type: GradeTest,
					Score: 81,
				},
				{
					Title: "Final Exam",
					Type: GradeQuiz,
					Score: 90,
				},

			},
		},
		{
			ID: 2,
			FirstName: "Jarvis",
			LastName: "Cs",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type: GradeQuiz,
					Score: 85,
				},
				{
					Title: "First test",
					Type: GradeTest,
					Score: 88,
				},
				{
					Title: "Final Exam",
					Type: GradeQuiz,
					Score: 78,
				},

			},
		},
		{
			ID: 3,
			FirstName: "Jack",
			LastName: "Nick",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type: GradeQuiz,
					Score: 82,
				},
				{
					Title: "First test",
					Type: GradeTest,
					Score: 81,
				},
				{
					Title: "Final Exam",
					Type: GradeQuiz,
					Score: 84,
				},

			},
		},
	}
}
