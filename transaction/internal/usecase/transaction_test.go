package usecase_test

import (
	"context"
	"testing"

	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	mock_publisher "github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/mocks"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/repository"
	mock_repo "github.com/ShmelJUJ/software-engineering/transaction/internal/repository/mocks"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	transactionID = "test-id"
	reason        = "test-reason"
)

func transactionHelper(t *testing.T) (*mock_logger.MockLogger, *mock_repo.MockTransactionRepo, *mock_publisher.MockTransactionPublisher) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := mock_logger.NewMockLogger(mockCtrl)
	repo := mock_repo.NewMockTransactionRepo(mockCtrl)
	publisher := mock_publisher.NewMockTransactionPublisher(mockCtrl)

	return l, repo, publisher
}

func TestGetTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		transactionID string
	}

	ctx := context.Background()

	transaction := &model.Transaction{
		ID: transactionID,
	}
	someErr := repository.NewGetTransactionError("test err", nil)

	testcases := []struct {
		name                string
		args                args
		mock                func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedTransaction *model.Transaction
		expectedErr         error
	}{
		{
			name: "Successfully get transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Get transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().GetTransaction(ctx, transactionID).Return(transaction, nil).Times(1)
			},
			expectedTransaction: transaction,
			expectedErr:         nil,
		},
		{
			name: "Failed to get transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Get transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().GetTransaction(ctx, transactionID).Return(nil, someErr).Times(1)
			},
			expectedTransaction: nil,
			expectedErr:         someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			actualTransaction, err := transactionUsecase.GetTransaction(
				testcase.args.ctx,
				testcase.args.transactionID,
			)
			assert.Equal(t, actualTransaction, testcase.expectedTransaction)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx         context.Context
		transaction *model.Transaction
	}

	ctx := context.Background()

	transaction := &model.Transaction{
		ID: transactionID,
	}

	someErr := repository.NewGetTransactionError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully create transaction",
			args: args{
				ctx:         ctx,
				transaction: transaction,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Create transaction usecase", map[string]interface{}{
					"transaction": transaction,
				})
				mtr.EXPECT().CreateTransaction(ctx, transaction).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to create transaction",
			args: args{
				ctx:         ctx,
				transaction: transaction,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Create transaction usecase", map[string]interface{}{
					"transaction": transaction,
				})
				mtr.EXPECT().CreateTransaction(ctx, transaction).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			err := transactionUsecase.CreateTransaction(
				testcase.args.ctx,
				testcase.args.transaction,
			)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestGetTransactionStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		transactionID string
	}

	ctx := context.Background()

	transactionStatus := model.Succeeded
	someErr := repository.NewGetTransactionStatusError("test err", nil)

	testcases := []struct {
		name                      string
		args                      args
		mock                      func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedTransactionStatus model.TransactionStatus
		expectedErr               error
	}{
		{
			name: "Successfully get transaction status",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Get transaction status usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().GetTransactionStatus(ctx, transactionID).Return(transactionStatus, nil).Times(1)
			},
			expectedTransactionStatus: transactionStatus,
			expectedErr:               nil,
		},
		{
			name: "Failed to get transaction status",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Get transaction status usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().GetTransactionStatus(ctx, transactionID).Return(model.Undefined, someErr).Times(1)
			},
			expectedTransactionStatus: model.Undefined,
			expectedErr:               someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			actualTransactionStatus, err := transactionUsecase.GetTransactionStatus(
				testcase.args.ctx,
				testcase.args.transactionID,
			)
			assert.Equal(t, actualTransactionStatus, testcase.expectedTransactionStatus)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestCancelTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		transactionID string
		reason        string
	}

	ctx := context.Background()

	someErr := repository.NewCancelTransactionError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully cancel transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				reason:        reason,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Cancel transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
					"reason":         reason,
				})
				mtr.EXPECT().CancelTransaction(ctx, transactionID, reason).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to cancel transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				reason:        reason,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Cancel transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
					"reason":         reason,
				})
				mtr.EXPECT().CancelTransaction(ctx, transactionID, reason).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			err := transactionUsecase.CancelTransaction(
				testcase.args.ctx,
				testcase.args.transactionID,
				testcase.args.reason,
			)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestAcceptTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		transactionID string
		sender        *model.TransactionUser
	}

	ctx := context.Background()

	transaction := &model.Transaction{
		ID: "test-transaction",
	}
	sender := &model.TransactionUser{
		ID: "test-user",
	}

	someErr := repository.NewAcceptTransactionError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo, *mock_publisher.MockTransactionPublisher)
		expectedErr error
	}{
		{
			name: "Successfully accept transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				sender:        sender,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo, mtp *mock_publisher.MockTransactionPublisher) {
				ml.EXPECT().Debug("Accept transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().AcceptTransaction(ctx, transactionID, sender).Return(nil).Times(1)
				mtr.EXPECT().GetTransaction(ctx, transactionID).Return(transaction, nil).Times(1)
				mtp.EXPECT().PublishProcessedTransaction(dto.FromTransactionModel(transaction)).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to accept transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				sender:        sender,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo, _ *mock_publisher.MockTransactionPublisher) {
				ml.EXPECT().Debug("Accept transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().AcceptTransaction(ctx, transactionID, sender).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
		{
			name: "Failed to get transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				sender:        sender,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo, _ *mock_publisher.MockTransactionPublisher) {
				ml.EXPECT().Debug("Accept transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().AcceptTransaction(ctx, transactionID, sender).Return(nil).Times(1)
				mtr.EXPECT().GetTransaction(ctx, transactionID).Return(nil, someErr).Times(1)
			},
			expectedErr: someErr,
		},
		{
			name: "Failed to publish succeeded transaction",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				sender:        sender,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo, mtp *mock_publisher.MockTransactionPublisher) {
				ml.EXPECT().Debug("Accept transaction usecase", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().AcceptTransaction(ctx, transactionID, sender).Return(nil).Times(1)
				mtr.EXPECT().GetTransaction(ctx, transactionID).Return(transaction, nil).Times(1)
				mtp.EXPECT().PublishProcessedTransaction(dto.FromTransactionModel(transaction)).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo, publisher)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			err := transactionUsecase.AcceptTransaction(
				testcase.args.ctx,
				testcase.args.transactionID,
				testcase.args.sender,
			)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx                context.Context
		updatedTransaction *model.Transaction
	}

	ctx := context.Background()

	transaction := &model.Transaction{
		ID: "test-transaction",
	}

	someErr := repository.NewUpdateTransactionError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully update transaction",
			args: args{
				ctx:                ctx,
				updatedTransaction: transaction,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Update transaction usecase", map[string]interface{}{
					"updated_transaction": transaction,
				})
				mtr.EXPECT().UpdateTransaction(ctx, transaction).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to update transaction",
			args: args{
				ctx:                ctx,
				updatedTransaction: transaction,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Update transaction usecase", map[string]interface{}{
					"updated_transaction": transaction,
				})
				mtr.EXPECT().UpdateTransaction(ctx, transaction).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			err := transactionUsecase.UpdateTransaction(
				testcase.args.ctx,
				testcase.args.updatedTransaction,
			)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}

func TestChangeTransactionStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx           context.Context
		transactionID string
		status        model.TransactionStatus
	}

	ctx := context.Background()

	status := model.Succeeded

	someErr := repository.NewChangeTransactionStatusError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully change transaction status",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				status:        status,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Change transaction status", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().ChangeTransactionStatus(ctx, transactionID, status).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to change transaction status",
			args: args{
				ctx:           ctx,
				transactionID: transactionID,
				status:        status,
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Change transaction status", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mtr.EXPECT().ChangeTransactionStatus(ctx, transactionID, status).Return(someErr).Times(1)
			},
			expectedErr: someErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			l, repo, publisher := transactionHelper(t)
			testcase.mock(l, repo)

			transactionUsecase := usecase.NewTransactionUsecase(repo, publisher, l)

			err := transactionUsecase.ChangeTransactionStatus(
				testcase.args.ctx,
				testcase.args.transactionID,
				testcase.args.status,
			)
			assert.Equal(t, err, testcase.expectedErr)
		})
	}
}
