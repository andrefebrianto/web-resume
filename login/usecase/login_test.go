package usecase_test

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"devoratio.dev/web-resume/config"
	"devoratio.dev/web-resume/internal/errorx"
	"devoratio.dev/web-resume/login/usecase"
	"devoratio.dev/web-resume/login/usecase/usecasemock"
	"devoratio.dev/web-resume/model"
)

var _ = Describe("Login with user credential", func() {
	var (
		mockController *gomock.Controller

		authenticationUsecaseMock *usecasemock.MockAuthenticationUsecase

		commonCtx             context.Context
		loginUsecase          *usecase.Login
		identifier            = "devoratio"
		ownerAccountStub      model.Owner
		errorInvalidParameter *errorx.Error
		appConfig             *config.Application
	)

	BeforeEach(func() {
		gofakeit.Seed(time.Now().UnixNano())
		mockController = gomock.NewController(GinkgoT())

		authenticationUsecaseMock = usecasemock.NewMockAuthenticationUsecase(mockController)
		appConfig = &config.Application{
			Authentication: config.Authentication{
				SigningKey: []byte("veryverysecretsigningkey"),
			},
		}

		loginUsecase = usecase.NewUsecase(authenticationUsecaseMock, appConfig)

		gofakeit.Struct(&ownerAccountStub)

		commonCtx = context.Background()
		errorInvalidParameter = errorx.New(errorx.TypeInvalidParameter, "username or email or password is invalid", nil)
	})

	AfterEach(func() {
		mockController.Finish()
	})

	When("the user send the correct combination of username or email and password", func() {
		It("tells the user that the username or email or password is invalid", func(ctx SpecContext) {
			password := "veryverysecurepassword"

			authenticationUsecaseMock.EXPECT().Authenticate(commonCtx, identifier, password).Return(&ownerAccountStub, nil)

			result, err := loginUsecase.Login(commonCtx, identifier, password)
			Expect(err).Should(BeNil())
			Expect(result).Should(ContainSubstring("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"))
		}, SpecTimeout(time.Second*2))
	})

	When("the user send the incorrect combination of username or email and password", func() {
		It("tells the user that the username or email or password is invalid", func(ctx SpecContext) {
			password := "twinkling"

			authenticationUsecaseMock.EXPECT().Authenticate(commonCtx, identifier, password).Return(nil, errorInvalidParameter)

			result, err := loginUsecase.Login(commonCtx, identifier, password)
			Expect(err.(*errorx.Error).Code).Should(Equal(errorInvalidParameter.Code))
			Expect(err.(*errorx.Error).Message).Should(Equal(errorInvalidParameter.Message))
			Expect(err.(*errorx.Error).Type).Should(Equal(errorInvalidParameter.Type))
			Expect(result).Should(Equal(""))
		}, SpecTimeout(time.Second*2))
	})
})
