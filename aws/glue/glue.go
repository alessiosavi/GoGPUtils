package glueutils

import (
	"context"
	"errors"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	"log"
	"sync"
)

var glueClient *glue.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		glueClient = glue.New(glue.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}

func StartWorkflow(workflowName string, params map[string]string) error {
	if stringutils.IsBlank(workflowName) {
		return errors.New("workflow is empty")
	}
	workflow, err := glueClient.GetWorkflow(context.Background(), &glue.GetWorkflowInput{Name: aws.String(workflowName), IncludeGraph: aws.Bool(false)})
	if err != nil {
		return err
	}

	if len(params) > 0 {
		for k, v := range params {
			workflow.Workflow.DefaultRunProperties[k] = v
		}
		if _, err = glueClient.UpdateWorkflow(context.Background(), &glue.UpdateWorkflowInput{Name: aws.String(workflowName), DefaultRunProperties: workflow.Workflow.DefaultRunProperties}); err != nil {
			return err
		}
	}

	_, err = glueClient.StartWorkflowRun(context.Background(), &glue.StartWorkflowRunInput{Name: aws.String(workflowName)})
	return err
}

func UpdateJob(jobname string) error {
	_job, err := glueClient.GetJob(context.Background(), &glue.GetJobInput{JobName: aws.String(jobname)})
	if err != nil {
		return err
	}
	job := _job.Job
	if _, err = glueClient.UpdateJob(context.Background(), &glue.UpdateJobInput{
		JobName: job.Name,
		JobUpdate: &types.JobUpdate{
			CodeGenConfigurationNodes: job.CodeGenConfigurationNodes,
			Command:                   job.Command,
			Connections:               job.Connections,
			DefaultArguments:          job.DefaultArguments,
			Description:               job.Description,
			ExecutionClass:            job.ExecutionClass,
			ExecutionProperty:         job.ExecutionProperty,
			GlueVersion:               aws.String("4.0"),
			LogUri:                    job.LogUri,
			MaxRetries:                job.MaxRetries,
			NonOverridableArguments:   job.NonOverridableArguments,
			NotificationProperty:      job.NotificationProperty,
			NumberOfWorkers:           aws.Int32(2),
			Role:                      job.Role,
			SecurityConfiguration:     job.SecurityConfiguration,

			//SourceControlDetails: &types.SourceControlDetails{
			//	AuthStrategy: job.SourceControlDetails.AuthStrategy,
			//	AuthToken:    job.SourceControlDetails.AuthToken,
			//	Branch:       aws.String("prod"),
			//	Folder:       job.SourceControlDetails.Folder,
			//	Owner:        job.SourceControlDetails.Owner,
			//	Provider:     job.SourceControlDetails.Provider,
			//	Repository:   job.SourceControlDetails.Repository,
			//},
			Timeout:    job.Timeout,
			WorkerType: types.WorkerTypeG1x,
		},
	}); err != nil {
		return err
	}
	return nil
}

func PushRepo(jobname string) error {
	job, err := GetJob("qa-update-monetary")
	if err != nil {
		return err
	}
	cfg := job.Job.SourceControlDetails
	if _, err = glueClient.UpdateJobFromSourceControl(context.Background(), &glue.UpdateJobFromSourceControlInput{
		AuthStrategy:    cfg.AuthStrategy,
		AuthToken:       cfg.AuthToken,
		BranchName:      aws.String("qa"),
		CommitId:        nil,
		Folder:          cfg.Folder,
		JobName:         aws.String(jobname),
		Provider:        cfg.Provider,
		RepositoryName:  job.Job.SourceControlDetails.Repository,
		RepositoryOwner: cfg.Owner,
	}); err != nil {
		return err
	}
	return nil
}

func ListWorkflows() ([]string, error) {
	workflows, err := glueClient.ListWorkflows(context.Background(), &glue.ListWorkflowsInput{})
	if err != nil {
		return nil, err
	}

	workflowNames := workflows.Workflows

	continuationToken := workflows.NextToken

	for continuationToken != nil {
		workflows, err = glueClient.ListWorkflows(context.Background(), &glue.ListWorkflowsInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		workflowNames = append(workflowNames, workflows.Workflows...)
		continuationToken = workflows.NextToken
	}
	return workflowNames, nil
}

func GetWorkflow(workflowName string) (*glue.GetWorkflowOutput, error) {
	return glueClient.GetWorkflow(context.Background(), &glue.GetWorkflowInput{Name: aws.String(workflowName), IncludeGraph: aws.Bool(true)})
}

func DeleteWorkflow(workflowName string, deep bool) (*glue.DeleteWorkflowOutput, []error) {
	var errs []error

	if !deep {
		workflow, err := glueClient.DeleteWorkflow(context.Background(), &glue.DeleteWorkflowInput{Name: aws.String(workflowName)})
		errs = append(errs, err)
		return workflow, errs
	}

	workflow, err := GetWorkflow(workflowName)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	for _, node := range workflow.Workflow.Graph.Nodes {
		if node.TriggerDetails != nil {
			for _, action := range node.TriggerDetails.Trigger.Actions {
				if action.JobName != nil {
					if _, err = DeleteJob(*action.JobName); err != nil {
						errs = append(errs, err)
					}
				}
				if action.CrawlerName != nil {
					if _, err = DeleteCrawlers(*action.CrawlerName); err != nil {
						errs = append(errs, err)
					}
				}
			}
			if _, err = DeleteTriggers(*node.TriggerDetails.Trigger.Name); err != nil {
				errs = append(errs, err)
			}
		}
		if node.JobDetails != nil {
			for _, job := range node.JobDetails.JobRuns {
				if _, err = DeleteJob(*job.JobName); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	deleteWorkflow, err := glueClient.DeleteWorkflow(context.Background(), &glue.DeleteWorkflowInput{Name: aws.String(workflowName)})
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return deleteWorkflow, errs
}

func ListCrawlers() ([]string, error) {
	crawlers, err := glueClient.ListCrawlers(context.Background(), &glue.ListCrawlersInput{NextToken: nil})
	if err != nil {
		return nil, err
	}

	crawlerNames := crawlers.CrawlerNames

	continuationToken := crawlers.NextToken

	for continuationToken != nil {
		crawlers, err = glueClient.ListCrawlers(context.Background(), &glue.ListCrawlersInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		crawlerNames = append(crawlerNames, crawlers.CrawlerNames...)
		continuationToken = crawlers.NextToken
	}
	return crawlerNames, nil
}

func GetCrawler(crawlerName string) (*glue.GetCrawlerOutput, error) {
	return glueClient.GetCrawler(context.Background(), &glue.GetCrawlerInput{Name: aws.String(crawlerName)})
}

func DeleteCrawlers(crawlerName string) (*glue.DeleteCrawlerOutput, error) {
	crawler, err := GetCrawler(crawlerName)
	if err != nil {
		return nil, err
	}

	for _, classifier := range crawler.Crawler.Classifiers {
		if _, err = glueClient.DeleteClassifier(context.Background(), &glue.DeleteClassifierInput{Name: aws.String(classifier)}); err != nil {
			return nil, err
		}
	}
	return glueClient.DeleteCrawler(context.Background(), &glue.DeleteCrawlerInput{Name: aws.String(crawlerName)})
}

func ListClassifiers() ([]string, error) {
	classifiers, err := glueClient.GetClassifiers(context.Background(), &glue.GetClassifiersInput{})
	if err != nil {
		return nil, err
	}

	var classifierNames []string

	for _, classifierName := range classifiers.Classifiers {
		if classifierName.XMLClassifier != nil {
			classifierNames = append(classifierNames, *classifierName.XMLClassifier.Name)
		}
		if classifierName.CsvClassifier != nil {
			classifierNames = append(classifierNames, *classifierName.CsvClassifier.Name)
		}
		if classifierName.GrokClassifier != nil {
			classifierNames = append(classifierNames, *classifierName.GrokClassifier.Name)
		}
		if classifierName.JsonClassifier != nil {
			classifierNames = append(classifierNames, *classifierName.JsonClassifier.Name)
		}
	}

	continuationToken := classifiers.NextToken

	for continuationToken != nil {
		classifiers, err = glueClient.GetClassifiers(context.Background(), &glue.GetClassifiersInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		for _, classifierName := range classifiers.Classifiers {
			if classifierName.XMLClassifier != nil {
				classifierNames = append(classifierNames, *classifierName.XMLClassifier.Name)
			}
			if classifierName.CsvClassifier != nil {
				classifierNames = append(classifierNames, *classifierName.CsvClassifier.Name)
			}
			if classifierName.GrokClassifier != nil {
				classifierNames = append(classifierNames, *classifierName.GrokClassifier.Name)
			}
			if classifierName.JsonClassifier != nil {
				classifierNames = append(classifierNames, *classifierName.JsonClassifier.Name)
			}
		}
		continuationToken = classifiers.NextToken
	}
	return classifierNames, nil
}

func GetClassifier(classifierName string) (*glue.GetClassifierOutput, error) {
	return glueClient.GetClassifier(context.Background(), &glue.GetClassifierInput{Name: aws.String(classifierName)})
}
func DeleteClassifier(classifierName string) (*glue.DeleteClassifierOutput, error) {
	return glueClient.DeleteClassifier(context.Background(), &glue.DeleteClassifierInput{Name: aws.String(classifierName)})
}

func ListTriggers() ([]string, error) {
	triggers, err := glueClient.ListTriggers(context.Background(), &glue.ListTriggersInput{})
	if err != nil {
		return nil, err
	}
	triggerNames := triggers.TriggerNames

	continuationToken := triggers.NextToken

	for continuationToken != nil {
		triggers, err = glueClient.ListTriggers(context.Background(), &glue.ListTriggersInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		triggerNames = append(triggerNames, triggers.TriggerNames...)
		continuationToken = triggers.NextToken
	}
	return triggerNames, nil
}

func GetTrigger(triggerName string) (*glue.GetTriggerOutput, error) {
	return glueClient.GetTrigger(context.Background(), &glue.GetTriggerInput{Name: aws.String(triggerName)})
}

func DeleteTriggers(triggerName string) (*glue.DeleteTriggerOutput, error) {
	return glueClient.DeleteTrigger(context.Background(), &glue.DeleteTriggerInput{Name: aws.String(triggerName)})
}

func ListJobs() ([]string, error) {
	jobs, err := glueClient.ListJobs(context.Background(), &glue.ListJobsInput{})
	if err != nil {
		return nil, err
	}

	jobNames := jobs.JobNames

	continuationToken := jobs.NextToken

	for continuationToken != nil {
		jobs, err = glueClient.ListJobs(context.Background(), &glue.ListJobsInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		jobNames = append(jobNames, jobs.JobNames...)
		continuationToken = jobs.NextToken
	}
	return jobNames, nil
}

func GetJob(jobName string) (*glue.GetJobOutput, error) {
	return glueClient.GetJob(context.Background(), &glue.GetJobInput{JobName: aws.String(jobName)})
}
func DeleteJob(jobName string) (*glue.DeleteJobOutput, error) {
	return glueClient.DeleteJob(context.Background(), &glue.DeleteJobInput{JobName: aws.String(jobName)})
}

func ListConnections() ([]string, error) {
	connections, err := glueClient.GetConnections(context.Background(), &glue.GetConnectionsInput{})
	if err != nil {
		return nil, err
	}

	var connectionNames []string
	for _, connectionName := range connections.ConnectionList {
		connectionNames = append(connectionNames, *connectionName.Name)
	}

	continuationToken := connections.NextToken

	for continuationToken != nil {
		connections, err = glueClient.GetConnections(context.Background(), &glue.GetConnectionsInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		for _, connectionName := range connections.ConnectionList {
			connectionNames = append(connectionNames, *connectionName.Name)
		}
		continuationToken = connections.NextToken
	}
	return connectionNames, nil
}

func GetConnection(connectionName string) (*glue.GetConnectionOutput, error) {
	return glueClient.GetConnection(context.Background(), &glue.GetConnectionInput{Name: aws.String(connectionName)})
}

func DeleteConnection(connectionName string) (*glue.DeleteConnectionOutput, error) {
	return glueClient.DeleteConnection(context.Background(), &glue.DeleteConnectionInput{ConnectionName: aws.String(connectionName)})
}

func ListDatabases() ([]string, error) {
	databases, err := glueClient.GetDatabases(context.Background(), &glue.GetDatabasesInput{})
	if err != nil {
		return nil, err
	}
	var databaseNames []string
	for _, databaseName := range databases.DatabaseList {
		databaseNames = append(databaseNames, *databaseName.Name)
	}

	continuationToken := databases.NextToken

	for continuationToken != nil {
		databases, err = glueClient.GetDatabases(context.Background(), &glue.GetDatabasesInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		for _, databaseName := range databases.DatabaseList {
			databaseNames = append(databaseNames, *databaseName.Name)
		}
		continuationToken = databases.NextToken
	}
	return databaseNames, nil
}

func GetDatabase(databaseName string) (*glue.GetDatabaseOutput, error) {
	return glueClient.GetDatabase(context.Background(), &glue.GetDatabaseInput{Name: aws.String(databaseName)})
}

func DeleteDatabase(databaseName string) (*glue.DeleteDatabaseOutput, error) {
	tables, err := ListTables(databaseName)
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		if _, err = DeleteTable(databaseName, table); err != nil {
			return nil, err
		}
	}
	return glueClient.DeleteDatabase(context.Background(), &glue.DeleteDatabaseInput{Name: aws.String(databaseName)})
}

func ListTables(databaseName string) ([]string, error) {
	tables, err := glueClient.GetTables(context.Background(), &glue.GetTablesInput{DatabaseName: aws.String(databaseName)})
	if err != nil {
		return nil, err
	}

	var tableNames []string
	for _, tableName := range tables.TableList {
		tableNames = append(tableNames, *tableName.Name)
	}

	continuationToken := tables.NextToken

	for continuationToken != nil {
		tables, err = glueClient.GetTables(context.Background(), &glue.GetTablesInput{
			NextToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}
		for _, tableName := range tables.TableList {
			tableNames = append(tableNames, *tableName.Name)
		}
		continuationToken = tables.NextToken
	}
	return tableNames, nil
}

func GetTable(databaseName, tableName string) (*glue.GetTableOutput, error) {
	return glueClient.GetTable(context.Background(), &glue.GetTableInput{DatabaseName: aws.String(databaseName), Name: aws.String(tableName)})
}

func DeleteTable(databaseName, tableName string) (*glue.DeleteTableOutput, error) {
	return glueClient.DeleteTable(context.Background(), &glue.DeleteTableInput{DatabaseName: aws.String(databaseName), Name: aws.String(tableName)})
}

func ListWorkflowExecution(wfName string) ([]types.WorkflowRun, error) {
	runs, err := glueClient.GetWorkflowRuns(context.Background(), &glue.GetWorkflowRunsInput{
		Name:         aws.String(wfName),
		IncludeGraph: aws.Bool(false),
		MaxResults:   aws.Int32(1000),
	})
	if err != nil {
		return nil, err
	}

	var res []types.WorkflowRun
	res = append(res, runs.Runs...)
	continuationToken := runs.NextToken

	for continuationToken != nil {
		runs, err = glueClient.GetWorkflowRuns(context.Background(), &glue.GetWorkflowRunsInput{
			Name:         aws.String(wfName),
			IncludeGraph: aws.Bool(false),
			MaxResults:   aws.Int32(1000),
			NextToken:    continuationToken,
		})
		if err != nil {
			return res, err
		}
		continuationToken = runs.NextToken
		res = append(res, runs.Runs...)
	}

	return res, err
}

func ResetAllJobBookmark() error {
	jobs, err := ListJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		log.Println("Removing " + job)
		bookmark, err := ResetJobBookmark(job)
		if err != nil {
			return err
		}
		log.Println(helper.MarshalIndent(bookmark))
	}
	return nil
}
func ResetJobBookmark(jobName string) (*glue.ResetJobBookmarkOutput, error) {
	return glueClient.ResetJobBookmark(context.Background(), &glue.ResetJobBookmarkInput{
		JobName: aws.String(jobName),
		//RunId:   nil,
	})
}
