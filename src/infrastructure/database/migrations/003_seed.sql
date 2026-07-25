-- 003_seed.sql — Realistic data for POS AI-First MVP (tienda de abarrotes)

INSERT OR IGNORE INTO usuarios (id, nombre, pin_hash, rol) VALUES
(1, 'Admin', '$2a$10$rxyqum0rxJ9htmlD5QYWE.9CY1XeKQfq7b4BT3RVF.e71GvccInIC', 'admin'),
(2, 'Maria Cajera', '$2a$10$hnm.vyyIRYJS3u/zENTeBOiuvS85cMGd1mlT8xw8kxyETf.UFOp1G', 'cajero');

INSERT OR IGNORE INTO categorias (id, nombre, descripcion) VALUES
(1, 'Bebidas', 'Refrescos, aguas, jugos'),
(2, 'Lacteos', 'Leche, yogurt, quesos'),
(3, 'Panaderia', 'Pan, tortillas, galletas'),
(4, 'Snacks', 'Papas, dulces, chocolates'),
(5, 'Limpieza', 'Jabones, detergentes'),
(6, 'Abarrotes', 'Arroz, frijol, aceite'),
(7, 'Frutas y Verduras', 'Productos frescos'),
(8, 'Carnes', 'Pollo, res, embutidos');

INSERT OR IGNORE INTO productos (id, nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad) VALUES
(1,  'Coca Cola 600ml',       'BEB-001', 1, 22.00, 15.00, 48, 12, 'unidad'),
(2,  'Agua Natural 1L',       'BEB-002', 1, 15.00, 8.00,  36, 12, 'unidad'),
(3,  'Jugo Del Valle 1L',     'BEB-003', 1, 28.00, 19.00, 20, 6,  'unidad'),
(4,  'Pepsi 600ml',           'BEB-004', 1, 22.00, 15.00, 30, 12, 'unidad'),
(5,  'Cerveza Corona 355ml',  'BEB-005', 1, 32.00, 22.00, 24, 12, 'unidad'),
(6,  'Leche Entera 1L',       'LAC-001', 2, 26.00, 20.00, 15, 6,  'unidad'),
(7,  'Yogurt Natural 1kg',    'LAC-002', 2, 38.00, 28.00, 10, 4,  'unidad'),
(8,  'Queso Oaxaca 250g',     'LAC-003', 2, 55.00, 40.00, 8,  3,  'unidad'),
(9,  'Pan Bimbo Grande',      'PAN-001', 3, 58.00, 42.00, 12, 4,  'unidad'),
(10, 'Tortillas 1kg',         'PAN-002', 3, 22.00, 16.00, 20, 8,  'paquete'),
(11, 'Galletas Marias 400g',  'PAN-003', 3, 28.00, 19.00, 15, 5,  'unidad'),
(12, 'Sabritas Original 45g', 'SNK-001', 4, 18.00, 12.00, 40, 15, 'unidad'),
(13, 'Doritos 62g',           'SNK-002', 4, 22.00, 15.00, 35, 12, 'unidad'),
(14, 'Chocolate Carlos V',    'SNK-003', 4, 12.00, 7.00,  50, 20, 'unidad'),
(15, 'Chicles Trident',       'SNK-004', 4, 15.00, 9.00,  30, 10, 'unidad'),
(16, 'Jabon Zote 400g',       'LIM-001', 5, 18.00, 12.00, 25, 8,  'unidad'),
(17, 'Cloro 1L',              'LIM-002', 5, 22.00, 14.00, 15, 5,  'unidad'),
(18, 'Fabuloso 1L',           'LIM-003', 5, 35.00, 24.00, 12, 4,  'unidad'),
(19, 'Arroz 1kg',             'ABR-001', 6, 28.00, 20.00, 20, 8,  'kg'),
(20, 'Frijol Negro 1kg',      'ABR-002', 6, 35.00, 25.00, 15, 6,  'kg'),
(21, 'Aceite 1L',             'ABR-003', 6, 42.00, 30.00, 18, 6,  'litro'),
(22, 'Pasta Spaguetti 200g',  'ABR-004', 6, 12.00, 7.00,  30, 10, 'unidad'),
(23, 'Atun en Lata',          'ABR-005', 6, 22.00, 15.00, 25, 8,  'unidad'),
(24, 'Platano',               'FRU-001', 7, 18.00, 10.00, 15, 5,  'kg'),
(25, 'Tomate',                'FRU-002', 7, 25.00, 15.00, 10, 4,  'kg'),
(26, 'Cebolla',               'FRU-003', 7, 20.00, 12.00, 12, 4,  'kg'),
(27, 'Limon',                 'FRU-004', 7, 35.00, 20.00, 8,  3,  'kg'),
(28, 'Pechuga de Pollo',      'CAR-001', 8, 95.00, 70.00, 5,  3,  'kg'),
(29, 'Jamon de Pavo 250g',    'CAR-002', 8, 45.00, 32.00, 8,  3,  'unidad'),
(30, 'Salchicha Paquete',     'CAR-003', 8, 38.00, 25.00, 10, 4,  'paquete');

INSERT OR IGNORE INTO clientes (id, nombre, telefono, direccion) VALUES
(1, 'Juan Perez',       '5551234567', 'Calle Reforma 123'),
(2, 'Maria Garcia',     '5559876543', 'Av. Juarez 456'),
(3, 'Carlos Lopez',     '5554567890', 'Calle Hidalgo 789'),
(4, 'Ana Martinez',     '5557891234', 'Av. Insurgentes 321'),
(5, 'Roberto Sanchez',  '5552345678', 'Calle Morelos 654'),
(6, 'Laura Torres',     '5553456789', 'Calle Madero 100'),
(7, 'Pedro Ramirez',    '5556781234', 'Av. Universidad 250'),
(8, 'Sofia Hernandez',  '5558901234', 'Calle Allende 75'),
(9, 'Miguel Angel Diaz','5551112233', 'Av. Chapultepec 500'),
(10, 'Gabriela Flores', '5554445566', 'Calle 5 de Mayo 80'),
(11, 'Fernando Castro',  '5557778899', 'Av. Revolucion 320'),
(12, 'Patricia Ruiz',   '5552223344', 'Calle Victoria 45');

INSERT OR IGNORE INTO ventas (id, usuario_id, cliente_id, total, metodo_pago, created_at) VALUES
(1, 1, 1, 88.00,  'efectivo',      datetime('now', '-6 days', 'localtime')),
(2, 1, NULL, 57.00,  'tarjeta',    datetime('now', '-5 days', 'localtime')),
(3, 2, 2, 145.00, 'efectivo',      datetime('now', '-4 days', 'localtime')),
(4, 1, 3, 72.00,  'transferencia', datetime('now', '-3 days', 'localtime')),
(5, 2, NULL, 198.00, 'efectivo',   datetime('now', '-2 days', 'localtime')),
(6, 1, 4, 55.00,  'tarjeta',      datetime('now', '-1 day', 'localtime')),
(7, 1, 1, 122.00, 'efectivo',     datetime('now', '-1 day', 'localtime')),
(8, 2, 5, 85.00,  'efectivo',     datetime('now', 'localtime')),
(9, 1, NULL, 44.00,  'efectivo',   datetime('now', 'localtime')),
(10, 1, 2, 210.00, 'tarjeta',     datetime('now', 'localtime')),
(11, 2, 6, 130.00, 'efectivo',    datetime('now', '-7 days', 'localtime')),
(12, 1, 7, 95.00,  'tarjeta',     datetime('now', '-7 days', 'localtime')),
(13, 2, 8, 62.00,  'efectivo',    datetime('now', '-5 days', 'localtime')),
(14, 1, 9, 180.00, 'transferencia', datetime('now', '-4 days', 'localtime')),
(15, 2, 10, 48.00, 'efectivo',    datetime('now', '-3 days', 'localtime')),
(16, 1, 11, 275.00, 'tarjeta',    datetime('now', '-2 days', 'localtime')),
(17, 2, 12, 92.00, 'efectivo',    datetime('now', '-1 day', 'localtime')),
(18, 1, 6, 156.00, 'efectivo',    datetime('now', 'localtime')),
(19, 2, 3, 320.00, 'tarjeta',     datetime('now', 'localtime')),
(20, 1, 9, 67.00,  'efectivo',    datetime('now', 'localtime'));

INSERT OR IGNORE INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES
(1, 1, 2, 22.00, 44.00), (1, 12, 2, 18.00, 36.00), (1, 14, 1, 12.00, 12.00),
(2, 6, 1, 26.00, 26.00), (2, 10, 1, 22.00, 22.00), (2, 15, 1, 15.00, 15.00),
(3, 28, 1, 95.00, 95.00), (3, 25, 2, 25.00, 50.00),
(4, 2, 2, 15.00, 30.00), (4, 21, 1, 42.00, 42.00),
(5, 5, 4, 32.00, 128.00), (5, 12, 2, 18.00, 36.00), (5, 13, 2, 22.00, 44.00),
(6, 8, 1, 55.00, 55.00),
(7, 1, 3, 22.00, 66.00), (7, 4, 2, 22.00, 44.00), (7, 14, 1, 12.00, 12.00),
(8, 19, 2, 28.00, 56.00), (8, 22, 3, 12.00, 36.00),
(9, 3, 1, 28.00, 28.00), (9, 11, 1, 28.00, 28.00),
(10, 28, 2, 95.00, 190.00), (10, 26, 1, 20.00, 20.00),
(11, 1, 3, 22.00, 66.00), (11, 5, 2, 32.00, 64.00),
(12, 28, 1, 95.00, 95.00),
(13, 22, 2, 12.00, 24.00), (13, 19, 1, 28.00, 28.00), (13, 14, 1, 12.00, 12.00),
(14, 5, 3, 32.00, 96.00), (14, 1, 2, 22.00, 44.00), (14, 21, 1, 42.00, 42.00),
(15, 12, 1, 18.00, 18.00), (15, 13, 1, 22.00, 22.00), (15, 15, 1, 15.00, 15.00),
(16, 28, 2, 95.00, 190.00), (16, 8, 1, 55.00, 55.00), (16, 19, 1, 28.00, 28.00),
(17, 6, 2, 26.00, 52.00), (17, 7, 1, 38.00, 38.00),
(18, 1, 4, 22.00, 88.00), (18, 4, 2, 22.00, 44.00), (18, 24, 1, 18.00, 18.00),
(19, 28, 3, 95.00, 285.00), (19, 29, 1, 45.00, 45.00),
(20, 2, 2, 15.00, 30.00), (20, 3, 1, 28.00, 28.00), (20, 15, 1, 15.00, 15.00);

INSERT OR IGNORE INTO configuracion (clave, valor) VALUES
('nombre_negocio', 'Mi Tiendita'),
('moneda', 'MXN'),
('iva_porcentaje', '16');

-- Movimientos de inventario (entradas de compra y ajustes)
INSERT OR IGNORE INTO inventario_movimientos (id, producto_id, tipo, cantidad, stock_resultante, referencia_tipo, motivo, usuario_id, created_at) VALUES
(1,  1,  'entrada', 24, 72, 'compra', 'Reabastecimiento semanal Coca Cola', 1, datetime('now', '-7 days', 'localtime')),
(2,  2,  'entrada', 24, 60, 'compra', 'Reabastecimiento agua', 1, datetime('now', '-7 days', 'localtime')),
(3,  5,  'entrada', 24, 48, 'compra', 'Compra cerveza Corona', 1, datetime('now', '-5 days', 'localtime')),
(4,  12, 'entrada', 30, 70, 'compra', 'Reposicion Sabritas', 2, datetime('now', '-5 days', 'localtime')),
(5,  28, 'entrada', 10, 15, 'compra', 'Compra pollo fresco', 1, datetime('now', '-4 days', 'localtime')),
(6,  19, 'entrada', 15, 35, 'compra', 'Reabastecimiento arroz', 1, datetime('now', '-3 days', 'localtime')),
(7,  6,  'entrada', 12, 27, 'compra', 'Compra leche', 2, datetime('now', '-3 days', 'localtime')),
(8,  14, 'ajuste',  5,  55, NULL,     'Ajuste por conteo fisico', 1, datetime('now', '-2 days', 'localtime')),
(9,  25, 'entrada', 8,  18, 'compra', 'Compra tomate fresco', 1, datetime('now', '-1 day', 'localtime')),
(10, 1,  'salida',  10, 62, 'venta',  'Ventas del dia', 1, datetime('now', 'localtime')),
(11, 28, 'salida',  5,  10, 'venta',  'Ventas pollo', 1, datetime('now', 'localtime')),
(12, 5,  'salida',  9,  39, 'venta',  'Ventas cerveza semana', 2, datetime('now', 'localtime'));
