package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
)

type OrganizationRepository interface {
	GetOrganizationFromTender(ctx context.Context, tenderID uuid.UUID) (uuid.UUID, error)
	GetUserUUID(ctx context.Context, name string) (uuid.UUID, error)
	CheckOrganizationResponsible(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (bool, error)
	GetOrganizationFromBid(ctx context.Context, bidID uuid.UUID) (uuid.UUID, error)
}

// organizationRepository структура репозитория для работы с организациями
type organizationRepository struct {
	sqlDB *sql.DB
	db    *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Ошибка инициализации SQL DB: %v", err))
	}
	return &organizationRepository{db: db, sqlDB: sqlDB}
}

// GetOrganizationFromTender возвращает ID организации по Id тендера
func (r *organizationRepository) GetOrganizationFromTender(ctx context.Context, tenderID uuid.UUID) (uuid.UUID, error) {

	query := `SELECT organization_id FROM tenders WHERE id = $1;`

	row := r.sqlDB.QueryRowContext(ctx, query, tenderID)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("тендер не найден: %v", err)
	}

	return id, nil
}

func (r *organizationRepository) GetOrganizationFromBid(ctx context.Context, bidID uuid.UUID) (uuid.UUID, error) {

	query := `SELECT organization_id FROM bids WHERE id = $1;`

	row := r.sqlDB.QueryRowContext(ctx, query, bidID)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("предложение не найдено: %v", err)
	}

	return id, nil
}

// GetUserUUID возвращает ID пользователя по его имени
func (r *organizationRepository) GetUserUUID(ctx context.Context, name string) (uuid.UUID, error) {

	query := `SELECT id
FROM employee
WHERE username = $1;`

	row := r.sqlDB.QueryRowContext(ctx, query, name)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		log.Printf("Ошибка создания тест запроса: %v", err)
		return uuid.Nil, fmt.Errorf("пользователь не найден: %v", err)
	}

	return id, nil
}

func (r *organizationRepository) CheckOrganizationResponsible(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (bool, error) {

	query := `SELECT 1 from organization_responsible where organization_id = $1 AND user_id = $2;`

	var exists int

	err := r.sqlDB.QueryRowContext(ctx, query, organizationID, userID).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil

	} else if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return true, nil
}
