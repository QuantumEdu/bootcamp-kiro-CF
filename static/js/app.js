// POS AI-First — Alpine.js Components

function salesCart() {
    return {
        items: [],
        search: '',
        metodoPago: 'efectivo',
        addItem(product) {
            const existing = this.items.find(i => i.id === product.id);
            if (existing) { existing.cantidad++; }
            else { this.items.push({ id: product.id, nombre: product.nombre, precio: product.precio, cantidad: 1 }); }
        },
        removeItem(index) { this.items.splice(index, 1); },
        incrementItem(index) { this.items[index].cantidad++; },
        decrementItem(index) { if (this.items[index].cantidad > 1) this.items[index].cantidad--; else this.removeItem(index); },
        total() { return this.items.reduce((sum, item) => sum + (item.precio * item.cantidad), 0); },
        async checkout() {
            if (this.items.length === 0) return;
            const payload = { items: this.items.map(i => ({ producto_id: i.id, cantidad: i.cantidad, precio_unitario: i.precio })), metodo_pago: this.metodoPago };
            try {
                const resp = await fetch('/api/ventas', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
                if (resp.ok) { this.items = []; htmx.trigger(document.body, 'ventaCreada'); }
                else { const err = await resp.json(); alert('Error: ' + (err.error || 'No se pudo registrar')); }
            } catch (e) { alert('Error: ' + e.message); }
        }
    };
}

function addToCart(id, nombre, precio) {
    const el = document.querySelector('[x-data]');
    if (el && el.__x) { el.__x.$data.addItem({ id, nombre, precio }); }
}
