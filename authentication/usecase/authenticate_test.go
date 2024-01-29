package usecase_test

import (
	"context"
	"time"

	"devoratio.dev/web-resume/authentication/usecase"
	"devoratio.dev/web-resume/authentication/usecase/repositorymock"
	"devoratio.dev/web-resume/internal/errorx"
	"devoratio.dev/web-resume/model"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticate user credential", Label("authentication"), func() {
	var (
		mockController *gomock.Controller

		authenticationRepoMock *repositorymock.MockAuthenticationRepository

		authenticateUsecase *usecase.Authentication

		identifier            = "devoratio"
		commonCtx             context.Context
		ownerAccountStub      model.OwnerAccount
		errorInvalidParameter *errorx.Error
	)

	BeforeEach(func() {
		gofakeit.Seed(time.Now().UnixNano())
		mockController = gomock.NewController(GinkgoT())

		authenticationRepoMock = repositorymock.NewMockAuthenticationRepository(mockController)

		authenticateUsecase = usecase.NewUsecase(authenticationRepoMock)

		gofakeit.Struct(&ownerAccountStub)
		ownerAccountStub.Password = "$2a$12$hWASkUwEkcS1CbsyRRwoBew5r7qwmXwH4YJyP.S149hghOg77UEQW"

		commonCtx = context.Background()

		errorInvalidParameter = errorx.New(errorx.TypeInvalidParameter, "username or email or password is invalid", nil)
	})

	AfterEach(func() {
		mockController.Finish()
	})

	When("the user send the correct combination of username or email and password", func() {
		Context("there is a problem with the database connection", func() {
			It("tells the user", func() {
				password := "twinkling"

				authenticationRepoMock.EXPECT().GetOwnerByUsernameOrEmail(commonCtx, identifier).Return(nil, errorx.ErrInternal)

				result, err := authenticateUsecase.Authenticate(commonCtx, identifier, password)

				Expect(err).Should(Equal(errorx.ErrInternal))
				Expect(result).Should(BeNil())
			})
		})

		Context("owner data retrieved successfully", func() {
			It("sends access token and refresh token to the user", func() {
				password := "veryverysecurepassword"

				authenticationRepoMock.EXPECT().GetOwnerByUsernameOrEmail(commonCtx, identifier).Return(&ownerAccountStub, nil)

				result, err := authenticateUsecase.Authenticate(commonCtx, identifier, password)
				Expect(result).Should(Equal(&(ownerAccountStub.Owner)))
				Expect(err).Should(BeNil())
			})
		})

	})

	When("the user send username or email that does not exist in the database", func() {
		It("tells the user username or email or password is invalid", func(ctx SpecContext) {
			password := "twinkling"

			authenticationRepoMock.EXPECT().GetOwnerByUsernameOrEmail(commonCtx, identifier).Return(nil, errorx.ErrNotFound)

			result, err := authenticateUsecase.Authenticate(commonCtx, identifier, password)
			Expect(err.(*errorx.Error).Code).Should(Equal(errorInvalidParameter.Code))
			Expect(err.(*errorx.Error).Message).Should(Equal(errorInvalidParameter.Message))
			Expect(err.(*errorx.Error).Type).Should(Equal(errorInvalidParameter.Type))
			Expect(result).Should(BeNil())
		}, SpecTimeout(time.Second*2))
	})

	When("the user send wrong password", func() {
		Context("there is a problem with the database connection", func() {
			It("tells the user", func() {
				password := "twinkling"

				authenticationRepoMock.EXPECT().GetOwnerByUsernameOrEmail(commonCtx, identifier).Return(nil, errorx.ErrInternal)

				result, err := authenticateUsecase.Authenticate(commonCtx, identifier, password)
				Expect(err).Should(Equal(errorx.ErrInternal))
				Expect(result).Should(BeNil())
			})
		})

		Context("owner data retrieved successfully", func() {
			It("tells the user that the username or email or password is invalid", func(ctx SpecContext) {
				password := "twinkling"

				authenticationRepoMock.EXPECT().GetOwnerByUsernameOrEmail(commonCtx, identifier).Return(&ownerAccountStub, nil)

				result, err := authenticateUsecase.Authenticate(commonCtx, identifier, password)
				Expect(err.(*errorx.Error).Code).Should(Equal(errorInvalidParameter.Code))
				Expect(err.(*errorx.Error).Message).Should(Equal(errorInvalidParameter.Message))
				Expect(err.(*errorx.Error).Type).Should(Equal(errorInvalidParameter.Type))
				Expect(result).Should(BeNil())
			}, SpecTimeout(time.Second*2))
		})
	})
})
