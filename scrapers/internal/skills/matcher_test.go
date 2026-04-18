package skills

import "testing"

func TestIsFresherJob(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		want        bool
	}{
		{
			name:        "accepts explicit fresher role",
			title:       "Graduate Engineer Trainee",
			description: "Great fit for fresh graduates with strong fundamentals in Java and SQL.",
			want:        true,
		},
		{
			name:        "accepts role with 0 to 1 years experience",
			title:       "Software Engineer",
			description: "Looking for candidates with 0 to 1 years of experience in backend development.",
			want:        true,
		},
		{
			name:        "rejects generic software engineer role without early career signals",
			title:       "Software Engineer",
			description: "Build product features with React and Node.js in a fast-moving team.",
			want:        false,
		},
		{
			name:        "rejects senior title even if software role",
			title:       "Senior Software Engineer",
			description: "Hands-on role building distributed systems.",
			want:        false,
		},
		{
			name:        "rejects higher experience requirement",
			title:       "Software Engineer",
			description: "Requires 3+ years of experience with AWS and microservices.",
			want:        false,
		},
		{
			name:        "rejects level ii role",
			title:       "Software Engineer II",
			description: "Own backend services and mentor teammates.",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFresherJob(tt.title, tt.description)
			if got != tt.want {
				t.Fatalf("IsFresherJob(%q, %q) = %v, want %v", tt.title, tt.description, got, tt.want)
			}
		})
	}
}
