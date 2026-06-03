import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getSupplier, createSupplier, updateSupplier } from "../api/suppliers";
import { useToast } from "../context/ToastContext";

export default function SupplierForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new" || id === "create";

  const [form, setForm] = useState({
    name: "", contact_person: "", phone: "", email: "",
    address: "", city: "", state: "", country: "",
    payment_terms: "", notes: "",
  });
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();

  useEffect(() => {
    if (!isNew && id) {
      getSupplier(id).then(s => setForm({
        name: s.name, contact_person: s.contact_person || "", phone: s.phone || "",
        email: s.email || "", address: s.address || "", city: s.city || "",
        state: s.state || "", country: s.country || "",
        payment_terms: s.payment_terms || "", notes: s.notes || "",
      })).catch(console.error);
    }
  }, [id, isNew]);

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (isNew) {
        await createSupplier(form);
        showToast("Supplier created successfully");
      } else {
        await updateSupplier(id!, form);
        showToast("Supplier updated successfully");
      }
      navigate("/suppliers");
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-page">
      <h1>{isNew ? "New Supplier" : "Edit Supplier"}</h1>
      <form onSubmit={handleSubmit}>
        <fieldset>
          <legend>Supplier Info</legend>
          <label>Name<input value={form.name} onChange={update("name")} required /></label>
          <label>Contact Person<input value={form.contact_person} onChange={update("contact_person")} /></label>
          <label>Phone<input value={form.phone} onChange={update("phone")} /></label>
          <label>Email<input type="email" value={form.email} onChange={update("email")} /></label>
        </fieldset>
        <fieldset>
          <legend>Address</legend>
          <label>Address<input value={form.address} onChange={update("address")} /></label>
          <label>City<input value={form.city} onChange={update("city")} /></label>
          <label>State<input value={form.state} onChange={update("state")} /></label>
          <label>Country<input value={form.country} onChange={update("country")} /></label>
        </fieldset>
        <fieldset>
          <legend>Payment & Notes</legend>
          <label>Payment Terms<input value={form.payment_terms} onChange={update("payment_terms")} placeholder="e.g. Net 30" /></label>
          <label>Notes<textarea value={form.notes} onChange={update("notes")} rows={3} /></label>
        </fieldset>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}
