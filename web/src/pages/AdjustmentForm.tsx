import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createAdjustment } from "../api/inventory";
import { listLocations, type Location } from "../api/locations";
import { listProducts, type Product } from "../api/products";

const MOVEMENT_TYPES = ["waste", "damage", "transfer_out", "theft", "correction"];

export default function AdjustmentForm() {
  const navigate = useNavigate();
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [form, setForm] = useState({
    product_id: "", location_id: "", batch_id: "", quantity: 1,
    movement_type: "waste", notes: "",
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

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await createAdjustment(form);
      navigate("/inventory");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <h1>Stock Adjustment</h1>
      <form onSubmit={handleSubmit}>
        <label>Product<select value={form.product_id} onChange={update("product_id")} required>
          <option value="">Select product...</option>
          {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
        </select></label>
        <label>Location<select value={form.location_id} onChange={update("location_id")} required>
          {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
        </select></label>
        <label>Type<select value={form.movement_type} onChange={update("movement_type")} required>
          {MOVEMENT_TYPES.map(t => <option key={t} value={t}>{t.replace("_", " ").replace(/\b\w/g, c => c.toUpperCase())}</option>)}
        </select></label>
        <label>Quantity<input type="number" min="1" value={form.quantity} onChange={e => setForm({ ...form, quantity: parseInt(e.target.value) || 0 })} required /></label>
        <label>Notes<textarea value={form.notes} onChange={update("notes")} rows={2} /></label>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Record Adjustment"}</button>
      </form>
    </div>
  );
}
