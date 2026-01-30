export interface Job {
  _id: string;
  externalId: string;
  title: string;
  company: string;
  location: string;
  description: string;
  url: string;
  source: string;
  postedAt: string;
  scrapedAt: string;
  remote: boolean;
  tags: string[];
  active: boolean;
}

export interface JobResponse {
  jobs: Job[];
  total: number;
  currentPage: number;
  totalPages: number;
}
