package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
)

type SQLiteDeviceRepository struct {
	db *sql.DB
}

func NewSQLiteDeviceRepository(db *sql.DB) *SQLiteDeviceRepository {
	return &SQLiteDeviceRepository{db: db}
}

func (r *SQLiteDeviceRepository) Save(device *domain.Device) error {
	_, err := r.db.Exec(
		`INSERT INTO devices (id, name, type, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		device.ID, device.Name, device.Type, device.Active,
		device.CreatedAt.UTC().Format(time.RFC3339Nano),
		device.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("failed to save device: %w", err)
	}
	return nil
}

func (r *SQLiteDeviceRepository) FindByID(id string) (*domain.Device, error) {
	row := r.db.QueryRow(
		`SELECT id, name, type, active, created_at, updated_at FROM devices WHERE id = ?`, id,
	)
	return r.scanDevice(row)
}

func (r *SQLiteDeviceRepository) Update(device *domain.Device) error {
	device.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE devices SET name = ?, type = ?, active = ?, updated_at = ? WHERE id = ?`,
		device.Name, device.Type, device.Active,
		device.UpdatedAt.UTC().Format(time.RFC3339Nano),
		device.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}
	return nil
}

func (r *SQLiteDeviceRepository) ListAll() ([]*domain.Device, error) {
	rows, err := r.db.Query(
		`SELECT id, name, type, active, created_at, updated_at FROM devices ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	return r.scanDevices(rows)
}

func (r *SQLiteDeviceRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM devices WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrDeviceNotFound
	}

	return nil
}

func (r *SQLiteDeviceRepository) scanDevice(row *sql.Row) (*domain.Device, error) {
	var device domain.Device
	var createdAt, updatedAt string

	err := row.Scan(&device.ID, &device.Name, &device.Type, &device.Active, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDeviceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan device: %w", err)
	}

	device.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	device.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAt)

	return &device, nil
}

func (r *SQLiteDeviceRepository) scanDevices(rows *sql.Rows) ([]*domain.Device, error) {
	var devices []*domain.Device
	for rows.Next() {
		var device domain.Device
		var createdAt, updatedAt string

		err := rows.Scan(&device.ID, &device.Name, &device.Type, &device.Active, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device row: %w", err)
		}

		device.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		device.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAt)
		devices = append(devices, &device)
	}

	return devices, nil
}
