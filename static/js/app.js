// POS AI-First — Alpine.js Components

// Sales Cart Component
function salesCart() {
    return {
        items: [],
        search: '',
        metodoPago: 'efectivo',

        addItem(product) {
            const existing = this.items.find(i => i.id === product.id);
            if (existing) {
                existing.cantidad++;
            } else {
                this.items.push({
                    id: product.id,
                    nombre: product.nombre,
                    precio: product.precio,
                    cantidad: 1
                });
            }
        },

        removeItem(index) {
            this.items.splice(index, 1);
        },

        incrementItem(index) {
            this.items[index].cantidad++;
        },

        decrementItem(index) {
            if (this.items[index].cantidad > 1) {
                this.items[index].cantidad--;
            } else {
                this.removeItem(index);
            }
        },

        total() {
            return this.items.reduce((sum, item) => sum + (item.precio * item.cantidad), 0);
        },

        async checkout() {
            if (this.items.length === 0) return;

            const payload = {
                items: this.items.map(i => ({
                    producto_id: i.id,
                    cantidad: i.cantidad,
                    precio_unitario: i.precio
                })),
                metodo_pago: this.metodoPago
            };

            try {
                const resp = await fetch('/api/ventas', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (resp.ok) {
                    this.items = [];
                    this.metodoPago = 'efectivo';
                    // Trigger HTMX refresh of recent sales
                    htmx.trigger(document.body, 'ventaCreada');
                    alert('Venta registrada correctamente!');
                } else {
                    const err = await resp.json();
                    alert('Error: ' + (err.error || 'No se pudo registrar la venta'));
                }
            } catch (e) {
                alert('Error de conexion: ' + e.message);
            }
        }
    };
}

// Global: Add product to cart from HTMX search results
function addToCart(id, nombre, precio) {
    const cartComponent = document.querySelector('[x-data]')?.__x;
    if (cartComponent) {
        cartComponent.$data.addItem({ id, nombre, precio });
    }
}
