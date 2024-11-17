package auth

import "github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/auth"

func MapDtoHashPasswordReqToPbHashPasswordReq(dto HashPasswordReq) *auth.HashPasswordReq {
	return &auth.HashPasswordReq{
		Password: dto.Password,
	}
}

func MapPbHashPasswordResToDtoHashPasswordRes(pb *auth.HashPasswordRes) HashPasswordRes {
	return HashPasswordRes{
		HashedPassword: pb.HashedPassword,
	}
}

func MapDtoVerifyTokenReqToPbVerifyTokenReq(dto VerifyTokenReq) *auth.VerifyTokenReq {
	return &auth.VerifyTokenReq{
		Token: dto.Token,
	}
}

func MapPbVerifyTokenResToDtoVerifyTokenRes(pb *auth.VerifyTokenRes) VerifyTokenRes {
	return VerifyTokenRes{
		ID:       int(pb.Id),
		Username: pb.Username,
		Email:    pb.Email,
		Role:     pb.Role,
	}
}
