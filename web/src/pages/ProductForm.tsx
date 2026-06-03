import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getProduct, createProduct, updateProduct, listCategories, barcodeLookup, type ProductCategory } from "../api/products";

const CLASSIFICATIONS = ["OTC", "Prescription", "Controlled", "Narcotic", "Device", "Supply"];
const FORMS = ["Tablet", "Capsule", "Syrup", "Suspension", "Injection", "Cream", "Ointment", "Drop", "Inhaler", "Spray", "Solution", "Powder"];
const UNITS = ["Tablet(s)", "Capsule(s)", "ml", "mg", "g", "mcg", "IU", "%", "Puff(s)", "Drop(s)"];

export default function ProductForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [form, setForm] = useState({
    name: "", description: "", category_id: "", classification: "",
    brand_name: "", generic_name: "", manufacturer: "", barcode: "", ndc: "",
    unit_of_measure: "", strength: "", form: "", reorder_level: 10,
  });
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    listCategories().then(setCategories).catch(console.error);
    if (!isNew && id && id !== "new") {
      getProduct(id).then(p => setForm({
        name: p.name, description: p.description || "",
        category_id: p.category_id || "", classification: p.classification || "",
        brand_name: p.brand_name || "", generic_name: p.generic_name || "",
        manufacturer: p.manufacturer || "", barcode: p.barcode || "", ndc: p.ndc || "",
        unit_of_measure: p.unit_of_measure || "", strength: p.strength || "",
        form: p.form || "", reorder_level: p.reorder_level || 10,
      })).catch(console.error);
    }
  }, [id, isNew]);

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleBarcodeLookup = async () => {
    if (!form.barcode) return;
    try {
      const p = await barcodeLookup(form.barcode);
      setForm({
        name: p.name, description: p.description || "",
        category_id: p.category_id || "", classification: p.classification || "",
        brand_name: p.brand_name || p.name, generic_name: p.generic_name || "",
        manufacturer: p.manufacturer || "", barcode: p.barcode || form.barcode,
        ndc: p.ndc || "", unit_of_measure: p.unit_of_measure || "",
        strength: p.strength || "", form: p.form || "",
        reorder_level: p.reorder_level || 10,
      });
    } catch {
      // barcode not found
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (isNew) {
        const product = await createProduct(form);
        navigate(`/products/${product.id}`);
      } else {
        await updateProduct(id!, form);
        navigate(`/products/${id}`);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-page">
      <h1>{isNew ? "New Product" : "Edit Product"}</h1>
      <form onSubmit={handleSubmit}>
        <fieldset>
          <legend>Basic Information</legend>
          <label>Name<input value={form.name} onChange={update("name")} required /></label>
          <label>Description<textarea value={form.description} onChange={update("description")} rows={2} /></label>
          <label>Category<select value={form.category_id} onChange={update("category_id")}>
            <option value="">Select...</option>
            {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select></label>
          <label>Classification<select value={form.classification} onChange={update("classification")}>
            <option value="">Select...</option>
            {CLASSIFICATIONS.map(c => <option key={c} value={c}>{c}</option>)}
          </select></label>
        </fieldset>

        <fieldset>
          <legend>Drug Details</legend>
          <label>Brand Name<input value={form.brand_name} onChange={update("brand_name")} /></label>
          <label>Generic Name<input value={form.generic_name} onChange={update("generic_name")} /></label>
          <label>Manufacturer<input value={form.manufacturer} onChange={update("manufacturer")} /></label>
          <label>Strength<input value={form.strength} onChange={update("strength")} placeholder="e.g. 500mg" /></label>
          <label>Form<select value={form.form} onChange={update("form")}>
            <option value="">Select...</option>
            {FORMS.map(f => <option key={f} value={f}>{f}</option>)}
          </select></label>
          <label>Unit of Measure<select value={form.unit_of_measure} onChange={update("unit_of_measure")}>
            <option value="">Select...</option>
            {UNITS.map(u => <option key={u} value={u}>{u}</option>)}
          </select></label>
        </fieldset>

        <fieldset>
          <legend>Barcode & NDC</legend>
          <label style={{ display: "flex", gap: 8, alignItems: "center" }}>
            Barcode
            <input value={form.barcode} onChange={update("barcode")} style={{ flex: 1 }} />
            <button type="button" onClick={handleBarcodeLookup} style={{ whiteSpace: "nowrap" }}>Lookup</button>
          </label>
          <label>NDC<input value={form.ndc} onChange={update("ndc")} /></label>
        </fieldset>

        <fieldset>
          <legend>Inventory Settings</legend>
          <label>Reorder Level<input type="number" value={form.reorder_level} onChange={e => setForm({ ...form, reorder_level: parseInt(e.target.value) || 0 })} /></label>
        </fieldset>

        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}
