-- 002_seed.sql
-- Seed data for POS AI-First MVP
-- Realistic data for a small PYME (convenience store / tienda de abarrotes)

-- ============================================================
-- USUARIOS (PIN: 1234 -> sha256 hash)
-- ============================================================
INSERT OR IGNORE INTO usuarios (id, nombre, pin_hash, rol) VALUES
(1, 'Admin', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 'admin'),
(2, 'Maria Cajera', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', 'cajero');
-- Admin PIN: 1234, Maria PIN: 123

-- ============================================================
-- CATEGORIAS
-- ============================================================
INSERT OR IGNORE INTO categorias (id, nombre, descripcion) VALUES
(1, 'Bebidas', 'Refrescos, aguas, jugos y bebidas alcoholicas'),
(2, 'Lacteos', 'Leche, yogurt, quesos'),
(3, 'Panaderia', 'Pan, tortillas, galletas'),
(4, 'Snacks', 'Papas, dulces, chocolates'),
(5, 'Limpieza', 'Jabones, detergentes, cloro'),
(6, 'Abarrotes', 'Arroz, frijol, aceite, pasta'),
(7, 'Frutas y Verduras', 'Productos frescos'),
(8, 'Carnes', 'Pollo, res, cerdo, embutidos');

-- ============================================================
-- PRODUCTOS
-- ============================================================
INSERT OR IGNORE INTO productos (id, nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad) VALUES
-- Bebidas
(1,  'Coca Cola 600ml',       'BEB-001', 1, 22.00, 15.00, 48, 12, 'unidad'),
(2,  'Agua Natural 1L',       'BEB-002', 1, 15.00, 8.00,  36, 12, 'unidad'),
(3,  'Jugo Del Valle 1L',     'BEB-003', 1, 28.00, 19.00, 20, 6,  'unidad'),
(4,  'Pepsi 600ml',           'BEB-004', 1, 22.00, 15.00, 30, 12, 'unidad'),
(5,  'Cerveza Corona 355ml',  'BEB-005', 1, 32.00, 22.00, 24, 12, 'unidad'),
-- Lacteos
(6,  'Leche Entera 1L',       'LAC-001', 2, 26.00, 20.00, 15, 6,  'unidad'),
(7,  'Yogurt Natural 1kg',    'LAC-002', 2, 38.00, 28.00, 10, 4,  'unidad'),
(8,  'Queso Oaxaca 250g',     'LAC-003', 2, 55.00, 40.00, 8,  3,  'unidad'),
-- Panaderia
(9,  'Pan Bimbo Grande',      'PAN-001', 3, 58.00, 42.00, 12, 4,  'unidad'),
(10, 'Tortillas 1kg',         'PAN-002', 3, 22.00, 16.00, 20, 8,  'paquete'),
(11, 'Galletas Marias 400g',  'PAN-003', 3, 28.00, 19.00, 15, 5,  'unidad'),
-- Snacks
(12, 'Sabritas Original 45g', 'SNK-001', 4, 18.00, 12.00, 40, 15, 'unidad'),
(13, 'Doritos 62g',           'SNK-002', 4, 22.00, 15.00, 35, 12, 'unidad'),
(14, 'Chocolate Carlos V',    'SNK-003', 4, 12.00, 7.00,  50, 20, 'unidad'),
(15, 'Chicles Trident',       'SNK-004', 4, 15.00, 9.00,  30, 10, 'unidad'),
-- Limpieza
(16, 'Jabon Zote 400g',       'LIM-001', 5, 18.00, 12.00, 25, 8,  'unidad'),
(17, 'Cloro 1L',              'LIM-002', 5, 22.00, 14.00, 15, 5,  'unidad'),
(18, 'Fabuloso 1L',           'LIM-003', 5, 35.00, 24.00, 12, 4,  'unidad'),
-- Abarrotes
(19, 'Arroz 1kg',             'ABR-001', 6, 28.00, 20.00, 20, 8,  'kg'),
(20, 'Frijol Negro 1kg',      'ABR-002', 6, 35.00, 25.00, 15, 6,  'kg'),
(21, 'Aceite 1L',             'ABR-003', 6, 42.00, 30.00, 18, 6,  'litro'),
(22, 'Pasta Spaguetti 200g',  'ABR-004', 6, 12.00, 7.00,  30, 10, 'unidad'),
(23, 'Atun en Lata',          'ABR-005', 6, 22.00, 15.00, 25, 8,  'unidad'),
-- Frutas y Verduras
(24, 'Platano',               'FRU-001', 7, 18.00, 10.00, 15, 5,  'kg'),
(25, 'Tomate',                'FRU-002', 7, 25.00, 15.00, 10, 4,  'kg'),
(26, 'Cebolla',               'FRU-003', 7, 20.00, 12.00, 12, 4,  'kg'),
(27, 'Limon',                 'FRU-004', 7, 35.00, 20.00, 8,  3,  'kg'),
-- Carnes
(28, 'Pechuga de Pollo',      'CAR-001', 8, 95.00, 70.00, 5,  3,  'kg'),
(29, 'Jamon de Pavo 250g',    'CAR-002', 8, 45.00, 32.00, 8,  3,  'unidad'),
(30, 'Salchicha Paquete',     'CAR-003', 8, 38.00, 25.00, 10, 4,  'paquete');

-- ============================================================
-- CLIENTES
-- ============================================================
INSERT OR IGNORE INTO clientes (id, nombre, telefono, direccion) VALUES
(1, 'Juan Perez',       '5551234567', 'Calle Reforma 123, Col. Centro'),
(2, 'Maria Garcia',     '5559876543', 'Av. Juarez 456, Col. Norte'),
(3, 'Carlos Lopez',     '5554567890', 'Calle Hidalgo 789'),
(4, 'Ana Martinez',     '5557891234', 'Av. Insurgentes 321'),
(5, 'Roberto Sanchez',  '5552345678', 'Calle Morelos 654');

-- ============================================================
-- VENTAS DE EJEMPLO (últimos días)
-- ============================================================
INSERT OR IGNORE INTO ventas (id, usuario_id, cliente_id, total, metodo_pago, created_at) VALUES
(1, 1, 1, 88.00,  'efectivo',      datetime('now', '-6 days', 'localtime')),
(2, 1, NULL, 57.00,  'tarjeta',    datetime('now', '-5 days', 'localtime')),
(3, 2, 2, 145.00, 'efectivo',      datetime('now', '-4 days', 'localtime')),
(4, 1, 3, 72.00,  'transferencia', datetime('now', '-3 days', 'localtime')),
(5, 2, NULL, 198.00, 'efectivo',    datetime('now', '-2 days', 'localtime')),
(6, 1, 4, 55.00,  'tarjeta',       datetime('now', '-1 day', 'localtime')),
(7, 1, 1, 122.00, 'efectivo',      datetime('now', '-1 day', 'localtime')),
(8, 2, 5, 85.00,  'mixto',         datetime('now', 'localtime')),
(9, 1, NULL, 44.00,  'efectivo',    datetime('now', 'localtime')),
(10, 1, 2, 210.00, 'tarjeta',      datetime('now', 'localtime'));

-- ============================================================
-- VENTA_ITEMS
-- ============================================================
INSERT OR IGNORE INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES
-- Venta 1
(1, 1, 2, 22.00, 44.00),
(1, 12, 2, 18.00, 36.00),
(1, 14, 1, 12.00, 12.00),
-- Venta 2
(2, 6, 1, 26.00, 26.00),
(2, 10, 1, 22.00, 22.00),
(2, 15, 1, 15.00, 15.00),
-- Venta 3
(3, 28, 1, 95.00, 95.00),
(3, 25, 2, 25.00, 50.00),
-- Venta 4
(4, 2, 2, 15.00, 30.00),
(4, 21, 1, 42.00, 42.00),
-- Venta 5
(5, 5, 4, 32.00, 128.00),
(5, 12, 2, 18.00, 36.00),
(5, 13, 2, 22.00, 44.00),
-- Venta 6
(6, 8, 1, 55.00, 55.00),
-- Venta 7
(7, 1, 3, 22.00, 66.00),
(7, 4, 2, 22.00, 44.00),
(7, 14, 1, 12.00, 12.00),
-- Venta 8
(8, 19, 2, 28.00, 56.00),
(8, 22, 3, 12.00, 36.00),
-- Venta 9
(9, 3, 1, 28.00, 28.00),
(9, 11, 1, 28.00, 28.00),
-- Venta 10
(10, 28, 2, 95.00, 190.00),
(10, 26, 1, 20.00, 20.00);

-- ============================================================
-- CONFIGURACION
-- ============================================================
INSERT OR IGNORE INTO configuracion (clave, valor) VALUES
('nombre_negocio', 'Mi Tiendita'),
('moneda', 'MXN'),
('iva_porcentaje', '16'),
('tickets_imprimir', '0');
