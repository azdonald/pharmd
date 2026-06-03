import { useEffect, useState } from "react";
import { listPrescribers, createPrescriber, updatePrescriber, deletePrescriber, type Prescriber } from "../api/prescribers";
import { useToast } from "../context/ToastContext";

export default function Prescribers() {
  const { showToast } = useToast();
  const [prescribers, setPrescribers] = useState<Prescriber[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState("");
  const [form, setForm] = useState({ name: "", license_number: "", phone: "", email: "", specialty: "", dea_number: "", npi_number: "", address: "" });

  const load = () => {
    setLoading(true);
    listPrescribers(page, 20, search).then(res => {
      setPrescribers(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, search]);

  const totalPages = Math.ceil(total / 20);

  const resetForm = () => {
    setShowForm(false);
    setEditId("");
    setForm({ name: "", license_number: "", phone: "", email: "", specialty: "", dea_number: "", npi_number: "", address: "" });
  };

  const openEdit = (p: Prescriber) => {
    setEditId(p.id);
    setForm({ name: p.name, license_number: p.license_number, phone: p.phone, email: p.email, specialty: p.specialty, dea_number: p.dea_number, npi_number: p.npi_number, address: p.address });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name) return;
    try {
      if (editId) {
        await updatePrescriber(editId, form);
      } else {
        await createPrescriber(form);
      }
      resetForm();
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this prescriber?")) return;
    try {
      await deletePrescriber(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Prescribers</h1>
        {!showForm && <button onClick={() => setShowForm(true)}>New Prescriber</button>}
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 16, padding: 16, border: "1px solid #ddd", borderRadius: 4 }}>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8, marginBottom: 8 }}>
            <div><label>Name *</label><input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required /></div>
            <div><label>License #</label><input value={form.license_number} onChange={e => setForm({ ...form, license_number: e.target.value })} /></div>
            <div><label>Phone</label><input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} /></div>
            <div><label>Email</label><input value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} /></div>
            <div><label>Specialty</label><input value={form.specialty} onChange={e => setForm({ ...form, specialty: e.target.value })} /></div>
            <div><label>DEA Number</label><input value={form.dea_number} onChange={e => setForm({ ...form, dea_number: e.target.value })} /></div>
            <div><label>NPI Number</label><input value={form.npi_number} onChange={e => setForm({ ...form, npi_number: e.target.value })} /></div>
            <div><label>Address</label><input value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} /></div>
          </div>
          <button type="submit">{editId ? "Update" : "Create"}</button>
          <button type="button" onClick={resetForm}>Cancel</button>
        </form>
      )}

      <input value={search} onChange={e => { setSearch(e.target.value); setPage(1); }} placeholder="Search prescribers..." style={{ width: "100%", marginBottom: 16 }} />

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Name</th><th>Specialty</th><th>License #</th><th>Phone</th><th>Active</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {prescribers.map(p => (
                <tr key={p.id}>
                  <td>{p.name}</td>
                  <td>{p.specialty || "—"}</td>
                  <td>{p.license_number || "—"}</td>
                  <td>{p.phone || "—"}</td>
                  <td>{p.is_active ? "Yes" : "No"}</td>
                  <td>
                    <button onClick={() => openEdit(p)}>Edit</button>
                    <button onClick={() => handleDelete(p.id)}>Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              {page > 1 && <button onClick={() => setPage(page - 1)}>Previous</button>}
              <span style={{ padding: "8px 0" }}>Page {page} of {totalPages}</span>
              {page < totalPages && <button onClick={() => setPage(page + 1)}>Next</button>}
            </div>
          )}
        </>
      )}
    </div>
  );
}
