import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createBatch } from "../api/inventory";
import { listLocations, type Location } from "../api/locations";
import { listProducts, type Product } from "../api/products";

export default function StockReceipt() {
  const navigate = useNavigate();
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [form, setForm] = useState({
    product_id: "", location_id: "", batch_number: "", quantity: 1,
    unit_cost: 0, selling_price: 0, manufacturing_date: "", expiry_date: "",
  });
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    Promise.all([
      listLocations().then(res => {
        setLocations(res.data);
        if (res.data.length > 0) setForm(f => ({ ...f, location_id: res.data[0].id }));
      }),
      listProducts(1, 500).then(res => setProducts(res.data)),
    ]).catch(console.error);
  }, []);

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await createBatch(form);
      navigate("/inventory");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <h1>Receive Stock</h1>
      <form onSubmit={handleSubmit}>
        <label>Product<select value={form.product_id} onChange={update("product_id")} required>
          <option value="">Select product...</option>
          {products.map(p => <option key={p.id} value={p.id}>{p.name} ({p.brand_name || p.generic_name || "—"})</option>)}
        </select></label>
        <label>Location<select value={form.location_id} onChange={update("location_id")} required>
          {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
        </select></label>
        <label>Batch Number<input value={form.batch_number} onChange={update("batch_number")} /></label>
        <label>Quantity<input type="number" min="1" value={form.quantity} onChange={e => setForm({ ...form, quantity: parseInt(e.target.value) || 0 })} required /></label>
        <label>Unit Cost<input type="number" step="0.01" min="0" value={form.unit_cost} onChange={e => setForm({ ...form, unit_cost: parseFloat(e.target.value) || 0 })} /></label>
        <label>Selling Price<input type="number" step="0.01" min="0" value={form.selling_price} onChange={e => setForm({ ...form, selling_price: parseFloat(e.target.value) || 0 })} /></label>
        <label>Manufacturing Date<input type="date" value={form.manufacturing_date} onChange={update("manufacturing_date")} /></label>
        <label>Expiry Date<input type="date" value={form.expiry_date} onChange={update("expiry_date")} /></label>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Receive Stock"}</button>
      </form>
    </div>
  );
}
