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
	"software engineer":       "Software Engineer",
	"software developer":      "Software Developer",
	"software":                "Software",
	"developer":               "Developer",
	"engineer":                "Engineer",
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

// IsFresherJob filters for fresher/intern/junior roles and excludes senior-level or advanced openings.
func IsFresherJob(title, description string) bool {
	text := strings.ToLower(strings.TrimSpace(title + " " + description))

	if hasAnyKeyword(text, seniorKeywords) || hasAnyKeyword(text, higherExperiencePatterns) {
		return false
	}

	juniorKeywords := []string{
		"intern", "internship", "junior", "entry level", "entry-level", "associate", "trainee", "fresher", "graduate", "campus",
		"apprentice", "new grad", "new graduate", "graduate program", "graduate engineer trainee",
		"sde 1", "sde i", "software engineer i", "software engineer 1", "developer i", "developer 1",
	}

	if hasAnyKeyword(text, juniorKeywords) {
		return true
	}

	experiencePatterns := []string{
		"0-1 year", "0-1 years", "0 to 1 year", "0 to 1 years",
		"0-2 year", "0-2 years", "0 to 2 years",
		"1 year", "1 years",
		"1-2 year", "1-2 years", "1 to 2 years",
		"less than 1 year", "less than 2 years",
		"up to 1 year", "up to 2 years",
		"no experience", "no prior experience", "fresh graduates", "recent graduates",
	}

	if hasAnyKeyword(text, experiencePatterns) {
		return true
	}

	return false
}

var seniorKeywords = []string{
	"senior", "sr", "sr.", "staff", "principal", "lead",
	"manager", "head", "director", "architect", "vp", "president", "supervisor",
	"mid senior", "mid-senior", "expert", "specialist", "consultant",
	"sde 2", "sde ii", "sde 3", "sde iii", "software engineer ii", "software engineer iii",
	"developer ii", "developer iii", "engineer ii", "engineer iii",
}

var higherExperiencePatterns = []string{
	"2+ year", "2+ years", "3+ year", "3+ years", "4+ year", "4+ years", "5+ year", "5+ years",
	"6+ year", "6+ years", "7+ year", "7+ years", "8+ year", "8+ years", "9+ year", "9+ years",
	"10+ year", "10+ years",
	"minimum 2 year", "minimum 2 years", "minimum 3 year", "minimum 3 years",
	"minimum 4 year", "minimum 4 years", "minimum 5 year", "minimum 5 years",
	"at least 2 year", "at least 2 years", "at least 3 year", "at least 3 years",
	"at least 4 year", "at least 4 years", "at least 5 year", "at least 5 years",
	"2 year experience", "2 years experience", "3 year experience", "3 years experience",
	"4 year experience", "4 years experience", "5 year experience", "5 years experience",
	"6 year experience", "6 years experience", "7 year experience", "7 years experience",
	"8 year experience", "8 years experience", "9 year experience", "9 years experience",
	"10 year experience", "10 years experience",
	"2 years of experience", "3 years of experience", "4 years of experience", "5 years of experience",
	"6 years of experience", "7 years of experience", "8 years of experience", "9 years of experience",
	"10 years of experience",
}

func hasAnyKeyword(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}
