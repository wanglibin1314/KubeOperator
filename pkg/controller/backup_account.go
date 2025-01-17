package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type BackupAccountController struct {
	Ctx                  context.Context
	BackupAccountService service.BackupAccountService
}

func NewBackupAccountController() *BackupAccountController {
	return &BackupAccountController{
		BackupAccountService: service.NewBackupAccountService(),
	}
}

// List BackupAccount
// @Tags backupAccounts
// @Summary Show all backupAccounts
// @Description Show backupAccounts
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /backupAccounts/ [get]
func (b BackupAccountController) Get() (*page.Page, error) {

	pg, _ := b.Ctx.Values().GetBool("page")
	if pg {
		num, _ := b.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := b.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return b.BackupAccountService.Page(num, size, condition.TODO())
	} else {
		var page page.Page
		projectName := b.Ctx.URLParam("projectName")
		items, err := b.BackupAccountService.List(projectName, condition.TODO())
		if err != nil {
			return &page, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

// Search BackupAccount
// @Tags backupAccounts
// @Summary Search backupAccount
// @Description Search backupAccount
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountRequest true "request"
// @Success 200 {object} dto.BackupAccount
// @Security ApiKeyAuth
// @Router /backupAccounts/search [post]
func (b BackupAccountController) PostSearch() (*page.Page, error) {
	var conditions condition.Conditions
	if b.Ctx.GetContentLength() > 0 {
		if err := b.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	p, _ := b.Ctx.Values().GetBool("page")
	if p {
		num, _ := b.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := b.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return b.BackupAccountService.Page(num, size, conditions)
	} else {
		var p page.Page
		projectName := b.Ctx.URLParam("projectName")
		items, err := b.BackupAccountService.List(projectName, condition.TODO())
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Create BackupAccount
// @Tags backupAccounts
// @Summary Create a backupAccount
// @Description create a backupAccount
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountRequest true "request"
// @Success 200 {object} dto.BackupAccount
// @Security ApiKeyAuth
// @Router /backupAccounts/ [post]
func (b BackupAccountController) Post() (*dto.BackupAccount, error) {
	var req dto.BackupAccountRequest
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	operator := b.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_BACKUP_ACCOUNT, req.Name)

	return b.BackupAccountService.Create(req)
}

func (b BackupAccountController) PostBatch() error {
	var req dto.BackupAccountOp
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = b.BackupAccountService.Batch(req)
	if err != nil {
		return err
	}

	operator := b.Ctx.Values().GetString("operator")
	delAccs := ""
	for _, item := range req.Items {
		delAccs += (item.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_BACKUP_ACCOUNT, delAccs)

	return err
}

// Update BackupAccount
// @Tags backupAccounts
// @Summary Update a backupAccount
// @Description Update a backupAccount
// @Accept  json
// @Produce  json
// @Param request body dto.BackupAccountRequest true "request"
// @Success 200 {object} dto.BackupAccount
// @Security ApiKeyAuth
// @Router /backupAccounts/{name}/ [patch]
func (b BackupAccountController) PatchBy(name string) (*dto.BackupAccount, error) {
	var req dto.BackupAccountUpdate
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	operator := b.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_BACKUP_ACCOUNT, name)

	return b.BackupAccountService.Update(name, req)
}

// Delete BackupAccount

// @Tags backupAccounts
// @Summary Delete a backupAccount
// @Description delete a  backupAccount by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /backupAccounts/{name}/ [delete]
func (b BackupAccountController) DeleteBy(name string) error {
	go kolog.Save("Delete", constant.DELETE_BACKUP_ACCOUNT, name)
	return b.BackupAccountService.Delete(name)
}

func (b BackupAccountController) PostBuckets() ([]interface{}, error) {
	var req dto.CloudStorageRequest
	err := b.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return b.BackupAccountService.GetBuckets(req)
}
