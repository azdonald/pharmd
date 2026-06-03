import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getSupplier, listSupplierProducts, setSupplierProducts, type Supplier, type SupplierProduct } from "../api/suppliers";
import { listProducts, type Product } from "../api/products";
import { useToast } from "../context/ToastContext";

export default function SupplierDetail() {
  const { id } = useParams();
  const [supplier, setSupplier] = useState<Supplier | null>(null);
  const [products, setProducts] = useState<SupplierProduct[]>([]);
  const [allProducts, setAllProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const { showToast } = useToast();
  const [newPrices, setNewPrices] = useState<Record<string, { product_id: string; unit_price: string; min_order_qty: string; lead_time_days: string; notes: string }>>({});

  const load = () => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      getSupplier(id),
      listSupplierProducts(id),
      listProducts(1, 500),
    ]).then(([s, p, all]) => {
      setSupplier(s);
      setProducts(p);
      setAllProducts(all.data);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const handleSavePrices = async () => {
    const items = Object.values(newPrices).filter(p => p.product_id && parseFloat(p.unit_price) > 0);
    if (items.length === 0) return;
    try {
      await setSupplierProducts(id!, items.map(p => ({
        product_id: p.product_id,
        unit_price: parseFloat(p.unit_price),
        min_order_qty: parseInt(p.min_order_qty) || 0,
        lead_time_days: parseInt(p.lead_time_days) || 0,
        notes: p.notes || undefined,
      })));
      setNewPrices({});
      setEditing(false);
      const updated = await listSupplierProducts(id!);
      setProducts(updated);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    }
  };

  const addPriceRow = () => {
    const key = Date.now().toString();
    setNewPrices({ ...newPrices, [key]: { product_id: "", unit_price: "", min_order_qty: "0", lead_time_days: "0", notes: "" } });
  };

  const updatePrice = (key: string, field: string, value: string) => {
    setNewPrices({ ...newPrices, [key]: { ...newPrices[key], [field]: value } });
  };

  if (loading) return <p>Loading...</p>;
  if (!supplier) return <p>Supplier not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>{supplier.name}</h1>
        <Link to={`/suppliers/${id}/edit`} className="btn">Edit</Link>
      </div>

      <table>
        <tbody>
          <tr><td><strong>Contact</strong></td><td>{supplier.contact_person || "—"}</td></tr>
          <tr><td><strong>Phone</strong></td><td>{supplier.phone || "—"}</td></tr>
          <tr><td><strong>Email</strong></td><td>{supplier.email || "—"}</td></tr>
          <tr><td><strong>Address</strong></td><td>{[supplier.address, supplier.city, supplier.state, supplier.country].filter(Boolean).join(", ") || "—"}</td></tr>
          <tr><td><strong>Payment Terms</strong></td><td>{supplier.payment_terms || "—"}</td></tr>
          <tr><td><strong>Notes</strong></td><td>{supplier.notes || "—"}</td></tr>
          <tr><td><strong>Status</strong></td><td>{supplier.is_active ? "Active" : "Inactive"}</td></tr>
        </tbody>
      </table>

      <h2 style={{ marginTop: 32 }}>Product Prices</h2>
      {products.length > 0 && (
        <table>
          <thead><tr><th>Product</th><th>Unit Price</th><th>Min Order</th><th>Lead Time</th><th>Notes</th></tr></thead>
          <tbody>
            {products.map(p => (
              <tr key={p.id}>
                <td>{p.product_name}</td>
                <td>{p.unit_price.toFixed(2)}</td>
                <td>{p.min_order_qty}</td>
                <td>{p.lead_time_days} days</td>
                <td>{p.notes || "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {!editing && <button onClick={() => setEditing(true)} style={{ marginTop: 12 }}>Set Prices</button>}

      {editing && (
        <div style={{ marginTop: 12 }}>
          {Object.entries(newPrices).map(([key, price]) => (
            <div key={key} style={{ display: "flex", gap: 8, marginBottom: 8, alignItems: "center" }}>
              <select value={price.product_id} onChange={e => updatePrice(key, "product_id", e.target.value)} required style={{ flex: 2 }}>
                <option value="">Select product...</option>
                {allProducts.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
              </select>
              <input type="number" step="0.01" value={price.unit_price} onChange={e => updatePrice(key, "unit_price", e.target.value)} placeholder="Price" style={{ flex: 1 }} required />
              <input type="number" value={price.min_order_qty} onChange={e => updatePrice(key, "min_order_qty", e.target.value)} placeholder="Min" style={{ width: 60 }} />
              <input type="number" value={price.lead_time_days} onChange={e => updatePrice(key, "lead_time_days", e.target.value)} placeholder="Days" style={{ width: 60 }} />
              <button onClick={() => { const { [key]: _, ...rest } = newPrices; setNewPrices(rest); }}>X</button>
            </div>
          ))}
          <div style={{ display: "flex", gap: 8 }}>
            <button onClick={addPriceRow}>Add Product</button>
            <button onClick={handleSavePrices}>Save Prices</button>
            <button onClick={() => { setEditing(false); setNewPrices({}); }}>Cancel</button>
          </div>
        </div>
      )}
    </div>
  );
}
