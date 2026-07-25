package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer conn.Close(ctx)

	// Seed users with bcrypt hashes (PIN 1234 and 123)
	_, err = conn.Exec(ctx, `
		INSERT INTO usuarios (id, nombre, pin_hash, rol) VALUES
		(1, 'Admin', '$2a$10$rxyqum0rxJ9htmlD5QYWE.9CY1XeKQfq7b4BT3RVF.e71GvccInIC', 'admin'),
		(2, 'Maria Cajera', '$2a$10$hnm.vyyIRYJS3u/zENTeBOiuvS85cMGd1mlT8xw8kxyETf.UFOp1G', 'cajero')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding users: %v", err)
	}

	// Seed categories
	_, err = conn.Exec(ctx, `
		INSERT INTO categorias (id, nombre, descripcion) VALUES
		(1, 'Bebidas', 'Refrescos, aguas, jugos'),
		(2, 'Lacteos', 'Leche, yogurt, quesos'),
		(3, 'Panaderia', 'Pan, tortillas, galletas'),
		(4, 'Snacks', 'Papas, dulces, chocolates'),
		(5, 'Limpieza', 'Jabones, detergentes'),
		(6, 'Abarrotes', 'Arroz, frijol, aceite')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding categories: %v", err)
	}

	// Seed products
	_, err = conn.Exec(ctx, `
		INSERT INTO productos (id, nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad) VALUES
		(1, 'Coca Cola 600ml', 'BEB-001', 1, 22.00, 15.00, 48, 12, 'unidad'),
		(2, 'Agua Natural 1L', 'BEB-002', 1, 15.00, 8.00, 36, 12, 'unidad'),
		(3, 'Jugo Del Valle 1L', 'BEB-003', 1, 28.00, 19.00, 20, 6, 'unidad'),
		(4, 'Leche Entera 1L', 'LAC-001', 2, 26.00, 20.00, 15, 6, 'unidad'),
		(5, 'Pan Bimbo Grande', 'PAN-001', 3, 58.00, 42.00, 12, 4, 'unidad'),
		(6, 'Sabritas Original 45g', 'SNK-001', 4, 18.00, 12.00, 40, 15, 'unidad'),
		(7, 'Arroz 1kg', 'ABR-001', 6, 28.00, 20.00, 20, 8, 'kg'),
		(8, 'Frijol Negro 1kg', 'ABR-002', 6, 35.00, 25.00, 15, 6, 'kg')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding products: %v", err)
	}

	// Seed clients
	_, err = conn.Exec(ctx, `
		INSERT INTO clientes (id, nombre, telefono, direccion) VALUES
		(1, 'Juan Perez', '5551234567', 'Calle Reforma 123'),
		(2, 'Maria Garcia', '5559876543', 'Av. Juarez 456'),
		(3, 'Carlos Lopez', '5554567890', 'Calle Hidalgo 789'),
		(4, 'Ana Martinez', '5557891234', 'Av. Insurgentes 321'),
		(5, 'Roberto Sanchez', '5552345678', 'Calle Morelos 654'),
		(6, 'Laura Torres', '5553456789', 'Calle Madero 100'),
		(7, 'Pedro Ramirez', '5556781234', 'Av. Universidad 250'),
		(8, 'Sofia Hernandez', '5558901234', 'Calle Allende 75'),
		(9, 'Miguel Angel Diaz', '5551112233', 'Av. Chapultepec 500'),
		(10, 'Gabriela Flores', '5554445566', 'Calle 5 de Mayo 80'),
		(11, 'Fernando Castro', '5557778899', 'Av. Revolucion 320'),
		(12, 'Patricia Ruiz', '5552223344', 'Calle Victoria 45')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding clients: %v", err)
	}

	// Seed some sales (spread across the last week)
	_, err = conn.Exec(ctx, `
		INSERT INTO ventas (id, usuario_id, cliente_id, total, metodo_pago, created_at) VALUES
		(1, 1, 1, 88.00, 'efectivo', NOW() - INTERVAL '6 days'),
		(2, 1, NULL, 57.00, 'tarjeta', NOW() - INTERVAL '5 days'),
		(3, 2, 2, 145.00, 'efectivo', NOW() - INTERVAL '4 days'),
		(4, 1, 3, 72.00, 'transferencia', NOW() - INTERVAL '3 days'),
		(5, 2, NULL, 198.00, 'efectivo', NOW() - INTERVAL '2 days'),
		(6, 1, 4, 55.00, 'tarjeta', NOW() - INTERVAL '1 day'),
		(7, 1, 1, 122.00, 'efectivo', NOW() - INTERVAL '1 day'),
		(8, 2, 5, 85.00, 'efectivo', NOW()),
		(9, 1, NULL, 44.00, 'efectivo', NOW()),
		(10, 1, 2, 210.00, 'tarjeta', NOW()),
		(11, 2, 6, 130.00, 'efectivo', NOW() - INTERVAL '7 days'),
		(12, 1, 7, 95.00, 'tarjeta', NOW() - INTERVAL '7 days'),
		(13, 2, 8, 62.00, 'efectivo', NOW() - INTERVAL '5 days'),
		(14, 1, 9, 180.00, 'transferencia', NOW() - INTERVAL '4 days'),
		(15, 2, 10, 48.00, 'efectivo', NOW() - INTERVAL '3 days'),
		(16, 1, 11, 275.00, 'tarjeta', NOW() - INTERVAL '2 days'),
		(17, 2, 12, 92.00, 'efectivo', NOW() - INTERVAL '1 day'),
		(18, 1, 6, 156.00, 'efectivo', NOW()),
		(19, 2, 3, 320.00, 'tarjeta', NOW()),
		(20, 1, 9, 67.00, 'efectivo', NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding sales: %v", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES
		(1, 1, 2, 22.00, 44.00), (1, 6, 2, 18.00, 36.00),
		(2, 4, 1, 26.00, 26.00), (2, 5, 1, 22.00, 22.00),
		(3, 7, 2, 28.00, 56.00), (3, 5, 1, 58.00, 58.00),
		(4, 2, 2, 15.00, 30.00), (4, 8, 1, 42.00, 42.00),
		(5, 5, 4, 32.00, 128.00), (5, 6, 2, 18.00, 36.00),
		(6, 8, 1, 55.00, 55.00),
		(7, 1, 3, 22.00, 66.00), (7, 4, 2, 22.00, 44.00),
		(8, 7, 2, 28.00, 56.00), (8, 3, 1, 12.00, 12.00),
		(9, 3, 1, 28.00, 28.00),
		(10, 1, 2, 95.00, 190.00), (10, 2, 1, 20.00, 20.00),
		(11, 1, 3, 22.00, 66.00), (11, 5, 2, 32.00, 64.00),
		(12, 1, 1, 95.00, 95.00),
		(13, 3, 2, 12.00, 24.00), (13, 7, 1, 28.00, 28.00),
		(14, 5, 3, 32.00, 96.00), (14, 1, 2, 22.00, 44.00),
		(15, 6, 1, 18.00, 18.00), (15, 6, 1, 22.00, 22.00),
		(16, 1, 2, 95.00, 190.00), (16, 8, 1, 55.00, 55.00),
		(17, 4, 2, 26.00, 52.00), (17, 7, 1, 38.00, 38.00),
		(18, 1, 4, 22.00, 88.00), (18, 4, 2, 22.00, 44.00),
		(19, 1, 3, 95.00, 285.00), (19, 2, 1, 45.00, 45.00),
		(20, 2, 2, 15.00, 30.00), (20, 3, 1, 28.00, 28.00)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding sale items: %v", err)
	}

	// Seed inventory movements
	_, err = conn.Exec(ctx, `
		INSERT INTO inventario_movimientos (id, producto_id, tipo, cantidad, stock_resultante, referencia_tipo, motivo, usuario_id, created_at) VALUES
		(1, 1, 'entrada', 24, 72, 'compra', 'Reabastecimiento semanal', 1, NOW() - INTERVAL '7 days'),
		(2, 2, 'entrada', 24, 60, 'compra', 'Reabastecimiento agua', 1, NOW() - INTERVAL '7 days'),
		(3, 5, 'entrada', 24, 48, 'compra', 'Compra cerveza', 1, NOW() - INTERVAL '5 days'),
		(4, 6, 'entrada', 30, 70, 'compra', 'Reposicion snacks', 2, NOW() - INTERVAL '5 days'),
		(5, 1, 'entrada', 10, 15, 'compra', 'Compra emergencia', 1, NOW() - INTERVAL '4 days'),
		(6, 7, 'entrada', 15, 35, 'compra', 'Reabastecimiento arroz', 1, NOW() - INTERVAL '3 days'),
		(7, 4, 'entrada', 12, 27, 'compra', 'Compra leche', 2, NOW() - INTERVAL '3 days'),
		(8, 1, 'salida', 10, 62, 'venta', 'Ventas del dia', 1, NOW()),
		(9, 5, 'salida', 9, 39, 'venta', 'Ventas cerveza', 2, NOW())
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding inventory movements: %v", err)
	}

	fmt.Println("✅ PostgreSQL seed completed successfully!")
}
