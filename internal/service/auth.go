package service

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/mall-api/global"
	"github.com/summerKK/mall-api/internal/dao"
	"github.com/summerKK/mall-api/internal/dto/admin"
	"github.com/summerKK/mall-api/internal/model"
	businessError "github.com/summerKK/mall-api/pkg/error"
	"github.com/summerKK/mall-api/pkg/security"
	"github.com/summerKK/mall-api/pkg/util"
)

type AdminService struct {
	service *service
	dao     *dao.Auth
}

func NewAdminService(ctx *gin.Context) *AdminService {
	svc := NewService(ctx)
	return &AdminService{
		service: svc,
		dao:     dao.NewAuth(global.DBEngine),
	}
}

// 用户登录
func (s *AdminService) Login(param *admin.UserLoginRequest) (user *model.UmsAdmin, err error) {
	defer func() {
		util.AddErrorToCtx(s.service.ctx, err)
	}()

	user, err = s.dao.GetItemByName(param.UserName)
	if err != nil {
		return
	}

	if security.VerifyPassword(user.Password, param.Password) {
		return
	}

	err = businessError.NewBusinessError("check auth failed")

	return
}

// 用户注册
func (s *AdminService) Register(param *admin.UserRegisterRequest) (user *model.UmsAdmin, err error) {
	defer func() {
		util.AddErrorToCtx(s.service.ctx, err)
	}()

	user = param.Convert2Model()
	// 查看用户是否已经存在
	existsUser, err := s.dao.GetItemByName(user.Username)
	if err != nil {
		return
	}

	if existsUser != nil {
		return nil, businessError.NewBusinessError("用户已存在")
	}

	hashedPassword, err := security.Hash(user.Password)
	if err != nil {
		return
	}

	user.Password = hashedPassword
	err = s.dao.Register(user)
	if err != nil {
		return
	}

	return
}

// 获取指定用户信息
func (s *AdminService) GetItem(userId uint) (user *model.UmsAdmin, err error) {
	defer func() {
		util.AddErrorToCtx(s.service.ctx, err)
	}()

	user = &model.UmsAdmin{}
	_, err = s.dao.GetItemById(userId, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 删除用户
func (s *AdminService) DeleteItem(userId int) (err error) {
	defer func() {
		util.AddErrorToCtx(s.service.ctx, err)
	}()

	err = s.dao.DeleteItemById(userId)

	return
}

// TODO 添加 menus、icon、roles
func (s *AdminService) GetAdminInfo() (userInfo map[string]interface{}, err error) {
	user, err := model.GetUser(s.service.ctx)
	if err != nil {
		return
	}

	userInfo = map[string]interface{}{
		"username": user.Username,
		"menus":    []int{},
		"icon":     user.Icon,
		"roles":    []int{},
	}

	return
}
