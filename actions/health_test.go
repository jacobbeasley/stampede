package actions

func (as *ActionSuite) Test_HealthCheck() {
	res := as.HTML("/api/health").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `"status":"healthy"`)
}

func (as *ActionSuite) Test_ReadyCheck() {
	res := as.HTML("/api/ready").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `"status":"ready"`)
}
