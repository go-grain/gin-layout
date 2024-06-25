package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/gin-layout/config"
	"github.com/go-grain/gin-layout/internal/repo/system/query"
	"github.com/go-grain/gin-layout/log"
	model "github.com/go-grain/gin-layout/model/system"
	"github.com/go-grain/gin-layout/pkg/encrypt"
	jwtx "github.com/go-grain/gin-layout/pkg/jwt"
	redisx "github.com/go-grain/gin-layout/pkg/redis"
	uuidx "github.com/go-grain/gin-layout/pkg/uuid"
)

type ISysUserRepo interface {
	Login(user *model.LoginReq) (*model.SysUser, error)
}

type SysUserService struct {
	repo ISysUserRepo
	rdb  redisx.IRedis
	conf *config.Config
	log  *log.Helper
}

func NewSysUserService(repo ISysUserRepo, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *SysUserService {
	return &SysUserService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  log.NewHelper(logger),
	}
}

func (s *SysUserService) InitSysUser() error {
	defaultAdminRole := s.conf.Role.DefaultAdminRole
	defaultRole := s.conf.Role.DefaultRole
	sysUser := []*model.SysUser{
		{UID: uuidx.UID(), Nickname: "哪个喵", Username: "admin", Password: encrypt.EncryptPassword("123456"), Roles: &model.Roles{defaultAdminRole, defaultRole}, Role: defaultAdminRole, Status: "yes"},
		{UID: uuidx.UID(), Nickname: "哪个喵", Username: "great", Password: encrypt.EncryptPassword("123456"), Roles: &model.Roles{defaultRole}, Role: defaultRole, Status: "yes"},
	}
	q := query.Q.SysUser
	count, err := q.Count()
	if err != nil {
		return err
	}
	// 有数据就默认已被初始化过,直接返回nil
	if count > 0 {
		return nil
	}

	return q.Create(sysUser...)
}

func (s *SysUserService) Login(login *model.LoginReq, ctx *gin.Context) (string, error) {
	user, err := s.repo.Login(login)
	if err != nil {
		return "", err
	}

	ctx.Set("LogType", "login")

	if !encrypt.ComparePasswords(user.Password, login.Password) {
		s.log.Errorw("errMsg", "用户登录", "err")
		return "", errors.New("账号或密码不正确")
	}

	if user.Status == "no" {
		s.log.Errorw("errMsg", "无法正常登录,账号已被冻结")
		return "", errors.New("无法正常登录,账号已被冻结")
	}

	jwt := jwtx.Jwt{}
	token, err := jwt.GenerateToken(user.UID, user.Role, s.conf.JWT.SecretKey, s.conf.JWT.ExpirationSeconds)
	if err != nil {
		s.log.Errorw("errMsg", "用户登录", "err", err.Error())
		return "", err
	}
	s.log.Infow("errMsg", "用户登录")
	return token, err
}
