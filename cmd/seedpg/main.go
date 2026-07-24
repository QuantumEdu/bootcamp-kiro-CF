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
		(2, 'Maria Garcia', '5559876543', 'Av. Juarez 456')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding clients: %v", err)
	}

	// Seed some sales
	_, err = conn.Exec(ctx, `
		INSERT INTO ventas (id, usuario_id, cliente_id, total, metodo_pago) VALUES
		(1, 1, 1, 88.00, 'efectivo'),
		(2, 1, NULL, 57.00, 'tarjeta'),
		(3, 2, 2, 145.00, 'efectivo')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding sales: %v", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES
		(1, 1, 2, 22.00, 44.00),
		(1, 6, 2, 18.00, 36.00),
		(2, 4, 1, 26.00, 26.00),
		(3, 7, 2, 28.00, 56.00),
		(3, 5, 1, 58.00, 58.00)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		log.Fatalf("seeding sale items: %v", err)
	}

	fmt.Println("✅ PostgreSQL seed completed successfully!")
}
