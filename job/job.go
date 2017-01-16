package job

type Job struct {
	job_title string
	job_company string
	job_description string
	job_published_date string
	job_compensation string
	job_location string
	isRemote string
	job_got_reference_from string
}

func New() Job {
	return &Job{}
}

func AddJob(newJob Job) {

}

func GetJob() Job {
	return  &Job{}
}
