package skills

import (
	"strings"
)

// SkillMap defines the normalization mapping for skills
var SkillMap = map[string]string{
	"javascript":              "JavaScript",
	"js":                      "JavaScript",
	"typescript":              "TypeScript",
	"ts":                      "TypeScript",
	"python":                  "Python",
	"java":                    "Java",
	"core java":               "Java",
	"j2ee":                    "Java EE",
	"golang":                  "Go",
	"go":                      "Go",
	"c++":                     "C++",
	"c#":                      "C#",
	"c sharp":                 "C#",
	"csharp":                  "C#",
	"dot net":                 ".NET",
	".net":                    ".NET",
	"dotnet":                  ".NET",
	"asp.net":                 "ASP.NET",
	"reactjs":                 "React",
	"react js":                "React",
	"react.js":                "React",
	"react":                   "React",
	"angular":                 "Angular",
	"angularjs":               "Angular",
	"vuejs":                   "Vue.js",
	"vue js":                  "Vue.js",
	"vue":                     "Vue.js",
	"nodejs":                  "Node.js",
	"node js":                 "Node.js",
	"node.js":                 "Node.js",
	"express":                 "Express.js",
	"expressjs":               "Express.js",
	"django":                  "Django",
	"flask":                   "Flask",
	"spring":                  "Spring",
	"spring boot":             "Spring Boot",
	"springboot":              "Spring Boot",
	"hibernate":               "Hibernate",
	"mysql":                   "MySQL",
	"postgresql":              "PostgreSQL",
	"postgres":                "PostgreSQL",
	"mongodb":                 "MongoDB",
	"mongo db":                "MongoDB",
	"sql server":              "SQL Server",
	"mssql":                   "SQL Server",
	"oracle":                  "Oracle",
	"redis":                   "Redis",
	"elasticsearch":           "Elasticsearch",
	"aws":                     "AWS",
	"amazon web services":     "AWS",
	"azure":                   "Azure",
	"gcp":                     "GCP",
	"google cloud":            "GCP",
	"docker":                  "Docker",
	"kubernetes":              "Kubernetes",
	"k8s":                     "Kubernetes",
	"terraform":               "Terraform",
	"jenkins":                 "Jenkins",
	"git":                     "Git",
	"github":                  "GitHub",
	"gitlab":                  "GitLab",
	"ci/cd":                   "CI/CD",
	"cicd":                    "CI/CD",
	"devops":                  "DevOps",
	"machine learning":        "Machine Learning",
	"ml":                      "Machine Learning",
	"deep learning":           "Deep Learning",
	"ai":                      "AI",
	"artificial intelligence": "AI",
	"data science":            "Data Science",
	"html":                    "HTML",
	"html5":                   "HTML5",
	"css":                     "CSS",
	"css3":                    "CSS3",
	"sass":                    "SASS",
	"scss":                    "SASS",
	"less":                    "LESS",
	"bootstrap":               "Bootstrap",
	"tailwind":                "Tailwind CSS",
	"jquery":                  "jQuery",
	"rest api":                "REST API",
	"restful":                 "REST API",
	"graphql":                 "GraphQL",
	"microservices":           "Microservices",
	"agile":                   "Agile",
	"scrum":                   "Scrum",
	"jira":                    "Jira",
}

// ExtractSkills returns a list of normalized skills found in the text
func ExtractSkills(text string) []string {
	text = strings.ToLower(text)
	found := make(map[string]bool)
	var output []string

	for key, normalized := range SkillMap {
		// Basic containment check.
		// Ideally we use word boundaries (regex), but this is faster for now.
		// "go" matches "good", so for short keys we need care.

		// Strict check for common English words that are also tech skills
		if key == "less" || key == "go" || key == "ai" {
			// Require case-insensitive match but surrounded by spaces/punctuation
			// or stricter context.
			// Ideally we use regex `\bkey\b`
			// For simplicity: check " key ", " key,", " key.", "/key", "(key)"
			strictKey := " " + key + " "
			if strings.Contains(text, strictKey) ||
				strings.Contains(text, " "+key+",") ||
				strings.Contains(text, "("+key+")") ||
				strings.Contains(text, "/"+key) {
				if !found[normalized] {
					output = append(output, normalized)
					found[normalized] = true
				}
			}
			continue
		}

		// Simple approach: Check if keyword exists
		if len(key) <= 3 {
			// For short keywords like 'go', 'js', 'c#', enforce some boundary-like checks manually or accept false positives for now
			// ' go ' or 'go,' etc.
			// Let's rely on simple Contains for now and assume job descriptions have context
			if strings.Contains(text, " "+key+" ") || strings.Contains(text, "/"+key) || strings.HasPrefix(text, key+" ") {
				if !found[normalized] {
					output = append(output, normalized)
					found[normalized] = true
				}
			}
		} else {
			if strings.Contains(text, key) {
				if !found[normalized] {
					output = append(output, normalized)
					found[normalized] = true
				}
			}
		}
	}
	return output
}

// IsSoftwareJob checks if the job has enough technical signals to be a software job
func IsSoftwareJob(text string) bool {
	matched := ExtractSkills(text)
	return len(matched) > 0
}
