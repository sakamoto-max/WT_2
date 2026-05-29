package fetcher
// mock the db
// create a jobs chan
// make it read data from the db and send it to the jobs chan

// consumer responsibilities :
// 1. listen to the result queue
// 2. convert the data into a job
// 3. push the job into the job queue

// consumer needs -> result queue, job queue

// producer responsibilities :
// 1. every time the ticker ticks -> fetch the data from the target outboxes
// 2. if data is found -> convert it into a job
// 3. push the job into the job queue

// producer needs -> dbs access & job queue 


// worker responsibilites :
// 1. listen to the jobs queue
// 2. when received a job -> perform the job
//

// worker needs -> jobs chan, queue & dbs access

// worker methods :
// 1. start -> starts listening to the jobs queue

// types of jobs :
// 1. pushing data into the target queues
// 2. updating the task db values

// job : 
// chan of type job
// methods :
// getTargetService
// getDbIndex
// getTaskName
// getNumberOfTries
// 

// 