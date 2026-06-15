import { useEffect, useState } from "react";
import { listPrescribers, createPrescriber, updatePrescriber, deletePrescriber, type Prescriber } from "../api/prescribers";
import { useToast } from "../context/ToastContext";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

const emptyForm = { name: "", license_number: "", phone: "", email: "", specialty: "", dea_number: "", npi_number: "", address: "" };

export default function Prescribers() {
  const { showToast } = useToast();
  const [prescribers, setPrescribers] = useState<Prescriber[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState("");
  const [form, setForm] = useState(emptyForm);
  const limit = 20;

  const load = () => {
    setLoading(true);
    listPrescribers(page, limit, search).then(res => {
      setPrescribers(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, search]);

  const totalPages = Math.ceil(total / limit);

  const resetForm = () => { setShowForm(false); setEditId(""); setForm(emptyForm); };

  const openEdit = (p: Prescriber) => {
    setEditId(p.id);
    setForm({ name: p.name, license_number: p.license_number, phone: p.phone, email: p.email, specialty: p.specialty, dea_number: p.dea_number, npi_number: p.npi_number, address: p.address });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name) return;
    try {
      if (editId) { await updatePrescriber(editId, form); }
      else { await createPrescriber(form); }
      resetForm(); load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this prescriber?")) return;
    try { await deletePrescriber(id); load(); }
    catch (err) { showToast(err instanceof Error ? err.message : "Delete failed", "error"); }
  };

  return (
    <div>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Prescribers</h2>
          <p className="text-body-lg text-on-surface-variant">Manage physicians and authorized prescribers</p>
        </div>
        {!showForm && (
          <button onClick={() => setShowForm(true)}
            className="btn-sky-action">
            <Icon name="add" className="mr-2" />New Prescriber
          </button>
        )}
      </div>

      {/* Inline form */}
      {showForm && (
        <form onSubmit={handleSubmit} className="mb-8 p-4 rounded-xl border border-outline-variant bg-surface-container-lowest">
          <h3 className="font-semibold text-on-surface mb-3">{editId ? "Edit Prescriber" : "New Prescriber"}</h3>
          <div className="grid grid-cols-2 gap-3 mb-3">
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Name *</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} required
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">License #</label>
              <input value={form.license_number} onChange={e => setForm({ ...form, license_number: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Phone</label>
              <input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Email</label>
              <input value={form.email} onChange={e => setForm({ ...form, email: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Specialty</label>
              <input value={form.specialty} onChange={e => setForm({ ...form, specialty: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">DEA Number</label>
              <input value={form.dea_number} onChange={e => setForm({ ...form, dea_number: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">NPI Number</label>
              <input value={form.npi_number} onChange={e => setForm({ ...form, npi_number: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Address</label>
              <input value={form.address} onChange={e => setForm({ ...form, address: e.target.value })}
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
          </div>
          <div className="flex gap-2">
            <button type="submit"
              className="px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container transition-all">{editId ? "Update" : "Create"}</button>
            <button type="button" onClick={resetForm}
              className="px-4 py-2 border border-outline-variant rounded-lg text-on-surface hover:bg-surface-container-high transition-all">Cancel</button>
          </div>
        </form>
      )}

      {/* Search */}
      <div className="mb-8 max-w-md">
        <div className="relative">
          <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant" />
          <input value={search} onChange={e => { setSearch(e.target.value); setPage(1); }}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest py-2 pl-10 pr-4 text-sm outline-none focus:ring-2 focus:ring-primary"
            placeholder="Search prescribers..." type="text" />
        </div>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading prescribers...</p></div>
          ) : prescribers.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No prescribers found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Specialty</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">License #</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Phone</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Active</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {prescribers.map(p => (
                  <tr key={p.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{p.name}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.specialty || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.license_number || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.phone || "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        p.is_active ? "bg-secondary-container/20 text-secondary" : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>{p.is_active ? "Active" : "Inactive"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <button onClick={() => openEdit(p)} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</button>
                        <button onClick={() => handleDelete(p.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} prescribers</p>
            <div className="flex space-x-2">
              <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_left" /></button>
              {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                const start = Math.max(1, Math.min(page - 2, totalPages - 4)); const p = start + i;
                if (p > totalPages) return null;
                return (
                  <button key={p} onClick={() => setPage(p)}
                    className={`w-10 h-10 rounded-md flex items-center justify-center font-bold text-sm ${
                      p === page ? "bg-primary text-on-primary" : "border border-outline-variant hover:bg-surface-container-high"
                    }`}>{p}</button>
                );
              })}
              <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_right" /></button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
