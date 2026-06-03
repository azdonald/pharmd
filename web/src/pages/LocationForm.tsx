import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getLocation, createLocation, updateLocation } from "../api/locations";
import { useToast } from "../context/ToastContext";

export default function LocationForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [form, setForm] = useState({ name: "", address: "", city: "", state: "", country: "", phone: "", email: "", tax_rate: 0, timezone: "UTC" });
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  useEffect(() => {
    if (!isNew && id) {
      getLocation(id).then(l => setForm({
        name: l.name, address: l.address || "", city: l.city || "", state: l.state || "",
        country: l.country || "", phone: l.phone || "", email: l.email || "",
        tax_rate: l.tax_rate || 0, timezone: l.timezone || "UTC",
      })).catch(console.error);
    }
  }, [id, isNew]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (isNew) {
        await createLocation(form);
        showToast("Location created successfully");
      } else {
        await updateLocation(id!, form);
        showToast("Location updated successfully");
      }
      navigate("/locations");
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-page">
      <h1>{isNew ? "New Location" : "Edit Location"}</h1>
      <form onSubmit={handleSubmit}>
        <label>Name<input value={form.name} onChange={update("name")} required /></label>
        <label>Address<input value={form.address} onChange={update("address")} /></label>
        <label>City<input value={form.city} onChange={update("city")} /></label>
        <label>State<input value={form.state} onChange={update("state")} /></label>
        <label>Country<input value={form.country} onChange={update("country")} /></label>
        <label>Phone<input value={form.phone} onChange={update("phone")} /></label>
        <label>Email<input type="email" value={form.email} onChange={update("email")} /></label>
        <label>Tax Rate<input type="number" step="0.01" value={form.tax_rate} onChange={e => setForm({ ...form, tax_rate: parseFloat(e.target.value) || 0 })} /></label>
        <label>Timezone<input value={form.timezone} onChange={update("timezone")} /></label>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}
