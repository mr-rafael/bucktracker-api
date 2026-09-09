package mapper

import (
	"github.com/Mr-Rafael/bucktracker-api/internal/dto"
	"github.com/Mr-Rafael/bucktracker-api/internal/service"
)

func ToLoginInput(reqParams dto.UserLoginRequestParams) service.LoginInput {
	return service.LoginInput{
		Email:    reqParams.Email,
		Password: reqParams.Password,
	}
}

func ToLoginResponse(loginInfo service.LoginInfo) dto.UserLoginResponseParams {
	return dto.UserLoginResponseParams{
		ID:          loginInfo.ID.String(),
		Email:       loginInfo.Email,
		Username:    loginInfo.UserName,
		AccessToken: loginInfo.AccessToken,
	}
}

func ToRefreshInput(refreshToken string) service.RefreshInput {
	return service.RefreshInput{
		RefreshToken: refreshToken,
	}
}

func ToRefreshResponse(refreshInfo service.RefreshInfo) dto.RefreshResponseParams {
	return dto.RefreshResponseParams{
		AccessToken: refreshInfo.AccessToken,
	}
}
