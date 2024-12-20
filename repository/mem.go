package repository

import (
	"clamp-core/models"
	"errors"
	"github.com/samber/lo"
	"log"
	"sync"

	"github.com/google/uuid"
)

type inMemoryRepository struct {
	workflows       map[string]*models.Workflow
	serviceRequests map[uuid.UUID]*models.ServiceRequest
	stepStatuses    map[uuid.UUID][]*models.StepsStatus
	lock            sync.RWMutex
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		workflows:       make(map[string]*models.Workflow),
		serviceRequests: make(map[uuid.UUID]*models.ServiceRequest),
		stepStatuses:    make(map[uuid.UUID][]*models.StepsStatus),
	}
}

func (repo *inMemoryRepository) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var results []*models.ServiceRequest
	for _, req := range repo.serviceRequests {
		if req.WorkflowName == workflowName {
			results = append(results, req)
		}
	}

	start := pageSize * pageNumber
	end := start + pageSize
	if start > len(results) {
		return []*models.ServiceRequest{}, nil
	}
	if end > len(results) {
		end = len(results)
	}
	return results[start:end], nil
}

func (repo *inMemoryRepository) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	statuses, exists := repo.stepStatuses[serviceRequestID]
	if !exists {
		return nil, nil
	}

	var results []*models.StepsStatus
	for _, status := range statuses {
		if status.StepID == stepID {
			results = append(results, status)
		}
	}
	return results, nil
}

func (repo *inMemoryRepository) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	statuses, exists := repo.stepStatuses[serviceRequestID]
	if !exists {
		return nil, nil
	}

	for _, s := range statuses {
		if s.StepID == stepID && s.Status == status {
			return s, nil
		}
	}
	return nil, errors.New("step status not found")
}

func (repo *inMemoryRepository) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	statuses, exists := repo.stepStatuses[serviceRequestID]
	if !exists {
		return nil, nil
	}

	var results []*models.StepsStatus
	for _, s := range statuses {
		if s.Status == status {
			results = append(results, s)
		}
	}
	return results, nil
}

func (repo *inMemoryRepository) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	statuses, exists := repo.stepStatuses[serviceRequestID]
	println()
	if !exists {
		return nil, nil
	}
	return statuses, nil
}

func (repo *inMemoryRepository) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	repo.stepStatuses[stepStatus.ServiceRequestID] = lo.UniqBy(append(repo.stepStatuses[stepStatus.ServiceRequestID], stepStatus), func(s *models.StepsStatus) int {
		return s.StepID
	})

	//val, _ := json.Marshal(repo.stepStatuses)
	//println(string(val))

	return stepStatus, nil
}

func (repo *inMemoryRepository) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	workflow, exists := repo.workflows[workflowName]
	if !exists {
		return nil, errors.New("workflow not found")
	}
	return workflow, nil
}

func (repo *inMemoryRepository) DeleteWorkflowByName(workflowName string) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	if _, exists := repo.workflows[workflowName]; !exists {
		return errors.New("workflow not found")
	}
	delete(repo.workflows, workflowName)
	return nil
}

func (repo *inMemoryRepository) SaveWorkflow(workflowReq *models.Workflow) (*models.Workflow, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	log.Println("------------------------------")
	log.Println(workflowReq.Name)
	repo.workflows[workflowReq.Name] = workflowReq
	return workflowReq, nil
}

func (repo *inMemoryRepository) FindServiceRequestByID(serviceRequestID uuid.UUID) (*models.ServiceRequest, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	req, exists := repo.serviceRequests[serviceRequestID]
	if !exists {
		return nil, errors.New("service request not found")
	}
	return req, nil
}

func (repo *inMemoryRepository) SaveServiceRequest(serviceReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	repo.serviceRequests[serviceReq.ID] = serviceReq
	return serviceReq, nil
}

func (repo *inMemoryRepository) GetWorkflows(pageNumber int, pageSize int, sortFields models.SortByFields) ([]*models.Workflow, int, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	var workflows []*models.Workflow
	for _, w := range repo.workflows {
		workflows = append(workflows, w)
	}

	// Sorting is omitted in this example, but it can be implemented if needed.
	totalWorkflowsCount := len(workflows)
	start := pageSize * (pageNumber - 1)
	end := start + pageSize
	if start > totalWorkflowsCount {
		return []*models.Workflow{}, 0, nil
	}
	if end > totalWorkflowsCount {
		end = totalWorkflowsCount
	}
	return workflows[start:end], totalWorkflowsCount, nil
}

func (repo *inMemoryRepository) Ping() error {
	// In-memory implementation always "pings" successfully.
	return nil
}
